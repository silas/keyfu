import os
import sys
import time
import unittest
import urllib
import bcrypt
import subprocess
import pymongo
from selenium import webdriver
from selenium.webdriver.common.keys import Keys as keys
from selenium.common import exceptions

# Import catalog
sys.path.insert(0, os.path.join('..', 'catalog'))
import loader

class Error(Exception): pass

class TestBasic(unittest.TestCase):

    URL='http://localhost:8000'
    EMAIL='silas@sewell.org'
    PASSWORD='pass123'

    def setUp(self):
        # setup mongodb connection
        self.conn = pymongo.Connection()
        self.conn.drop_database('keyfu')
        self.db = self.conn.keyfu

        # Setup test user
        self.db.user.insert({
            'email': self.EMAIL,
            'password': bcrypt.hashpw(self.PASSWORD, bcrypt.gensalt(12)),
        }, safe=True)

        # Load catalog
        loader.main(db=self.db)

        # Start browser
        self.d = webdriver.Firefox()

    def tearDown(self):
        self.d.close()

    @property
    def url(self):
        return self.d.current_url

    def wait(self, time=5):
        self.d.implicitly_wait(time)

    def home(self):
        if self.url != self.URL:
            self.d.get(self.URL)
            self.wait()
        return self.d.find_element_by_name('q')

    def login(self):
        if not self.d.get_cookie('sid'):
            if '/login' not in self.url:
                append = ''
                if self.url.startswith('http'):
                    append = '?' + urllib.urlencode({'url': self.url})
                self.d.get(self.URL + '/login' + append)
            self.login_form()
            self.wait()

    def login_form(self, email=None, password=None):
        self.assertEqual(self.d.title, 'Login - KeyFu')

        email = email or self.EMAIL
        password = password or self.PASSWORD

        e = self.d.find_element_by_name('email')
        e.send_keys(email)

        p = self.d.find_element_by_name('password')
        p.send_keys(password)

        p.submit()

    def test_auth(self):
        self.d.delete_all_cookies()

        # Prompt for login
        q = self.home()
        q.send_keys(':edit :')
        q.submit()

        self.wait()

        # Ensure we're at login page
        self.assertEqual(self.d.title, 'Login - KeyFu')

        def check_failure():
            self.assertEqual(self.d.title, 'Login - KeyFu')
            error = self.d.find_element_by_css_selector('div.error')
            self.assertEqual(error.text, 'hide\nInvalid email or password.')

        # Bad email
        self.login_form(email='fail')
        self.wait()
        check_failure()

        # Bad password
        self.login_form(password='fail')
        self.wait()
        check_failure()

        # Login successfully
        self.login()
        self.assertEqual(self.d.title, 'Add - KeyFu')

        # Logout succesfully
        self.d.get(self.URL + '/logout')
        self.assertEqual(self.url, self.URL + '/')

        logout = self.d.find_elements_by_css_selector('ul.util a')[2]
        self.assertEqual(logout.text, 'Login')

    def test_builtin(self):
        self.login()
        q = self.home()

        # Edit test keyword
        q.send_keys(':e example')
        q.submit()

        # Create keyword
        self.wait()
        t = self.d.find_element_by_name('type')
        t.send_keys('b')
        b = self.d.find_element_by_name('body')
        b.send_keys('com.keyfu.edit')
        b.submit()

        # Test keyword without query
        q = self.home()
        q.send_keys('example')
        q.submit()

        self.wait()
        self.assertEqual(self.url, self.URL + '/edit')

        # Test keyword with query
        q = self.home()
        q.send_keys('example test')
        q.submit()

        self.wait()
        self.assertEqual(
            self.url,
            self.URL + '/edit?' + urllib.urlencode({'q': 'test'})
        )

    def test_link(self):
        self.login()
        q = self.home()

        # Edit test keyword
        q.send_keys(':e example')
        q.submit()

        # Create keyword
        self.wait()
        b = self.d.find_element_by_name('body')
        b.send_keys('%s/test\n%s/test#test' % (self.URL, self.URL))
        b.submit()

        # Test keyword without query
        q = self.home()
        q.send_keys('example')
        q.submit()

        self.wait()
        self.assertEqual(self.url, '%s/test' % self.URL)

        # Test keyword with query
        q = self.home()
        q.send_keys('example test')
        q.submit()

        self.wait()
        self.assertEqual(self.url, '%s/test#test' % self.URL)

if __name__ == '__main__':
    unittest.main()
