from datetime import datetime,timedelta
import hashlib
import json
import logging
import os
import webapp2

from google.appengine.api import users
from google.appengine.ext import ndb

class PriceStore(ndb.Model):
  creationTimestamp = ndb.DateTimeProperty(auto_now_add=True)
  data = ndb.StringProperty()

  @staticmethod
  def ancestorKey(user):
    return ndb.Key("UserEmailHash", hashlib.sha256(user.email()).hexdigest())

  @staticmethod
  def save(user, data):
    PriceStore(parent = PriceStore.ancestorKey(user), data = data).put()

  @classmethod
  def load(cls, user):
    # Every call to save creates a new copy so make sure to get the latest one.
    data = cls.query(ancestor = PriceStore.ancestorKey(user)).order(-cls.creationTimestamp).fetch(1)
    if len(data):
      return data[0].data
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
