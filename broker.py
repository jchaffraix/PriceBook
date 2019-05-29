from datetime import datetime,timedelta
import hashlib
import json
import logging
import os
import webapp2

from google.appengine.api import users
from google.appengine.ext import ndb

OBJECT_KEY = u"obj_key"

class PriceStore(ndb.Model):
  creationTimestamp = ndb.DateTimeProperty(auto_now_add=True)
  # TODO: Validate the data to avoid persistent XSS.
  data = ndb.StringProperty()
  obj_key = ndb.StringProperty()

  @staticmethod
  def ancestorKey(user):
    return ndb.Key("UserEmailHash", hashlib.sha256(user.email()).hexdigest())

  @staticmethod
  def save(user, data):
    # This is super inefficient. This is because we have to fetch *all*
    # the keys to find the matching ones and reuse them.
    # TODO: Figure out how to do this efficiently (store the key in the object?).
    objects = json.JSONDecoder().decode(data)
    for object in objects:
      products = PriceStore.query(ancestor=PriceStore.ancestorKey(user))
      foundProduct = False
      for product in products.fetch():
        if (product.obj_key == object[OBJECT_KEY]):
          product.data = json.dumps(object)
          product.put()
          foundProduct = True
          break
      if not foundProduct:
        PriceStore(parent = PriceStore.ancestorKey(user), obj_key = object[u"obj_key"], data = json.dumps(object)).put()


  # TODO: We need a delete method.
  # TODO: This schema needs to be aligned with what we want to do.
  # TODO: Implement a way to fetch the key and just use the storage key for efficiency.

  @classmethod
  def load(cls, user):
    data = cls.query(ancestor = PriceStore.ancestorKey(user)).fetch()
    if len(data):
      # TODO: This is not great. Refactor this into a better abstraction.
      s = "["
      for d in data:
        s += d.data
        s += ","
      s = s[:-1]
      s += "]"
      return s
    # Return an empty JSON object as the UI expects JSON.
    return "[]"

  @classmethod
  def removeOldEntries(cls):
    deadline = datetime.now() - timedelta(days=60)
    logging.info("Starting cleanup with deadline = %s" % deadline)
    keys_to_delete = []
    for data in cls.query().iter():
      logging.info("Seen entity with creationTimestamp = %s" % data.creationTimestamp)
      if (data.creationTimestamp < deadline):
        logging.info("Removing entity with creationTimestamp = %s" % data.creationTimestamp)
        keys_to_delete.append(data.key())
    ndb.delete_multi(keys_to_delete)


class PriceBookPage(webapp2.RequestHandler):
  def get(self):
    html = open('main.html').read()
    self.response.write(html)


class PriceBookStorePage(webapp2.RequestHandler):
  def get(self):
    user = users.get_current_user()
    if not user:
      self.response.write("Need to log in to store information")
      self.response.status_int = 401
      return

    self.response.write(PriceStore.load(user))

  def post(self):
    user = users.get_current_user()
    if not user:
      self.response.write("Need to log in to store information")
      self.response.status_int = 401
      return

    logging.info(self.request.body)
    PriceStore.save(user, self.request.body)

# This works as we use login: required in app.yaml
# and that the window was created with window.open.
class LoginPage(webapp2.RequestHandler):
  def get(self):
    self.response.write("<!DOCTYPE html><script>window.close();</script>");

class RemoveOldEntriesPage(webapp2.RequestHandler):
  def get(self):
    PriceStore.removeOldEntries()
    self.response.write("Done")

app = webapp2.WSGIApplication([
    ('/login', LoginPage),
    ('/store', PriceBookStorePage),
    ('/tasks/removeOldEntries', RemoveOldEntriesPage),
    ('/', PriceBookPage),
# TODO: Should I switch to debug=False?
], debug=True)
