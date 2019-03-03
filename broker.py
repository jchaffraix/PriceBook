import hashlib
import json
import logging
import os
import webapp2

from google.appengine.api import users
from google.appengine.ext import ndb

class PriceBookPage(webapp2.RequestHandler):
  def get(self):
    html = open('main.html').read()
    self.response.write(html)

app = webapp2.WSGIApplication([
    ('/', PriceBookPage),
# TODO: Should I switch to debug=False?
], debug=True)
