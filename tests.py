# coding: utf-8

import json
import unittest

import requests


with open('configs/testing.json') as f:
    config = json.loads(f.read())

base_url = 'http://%s:%d' % (config['Host'], config['Port'])

test_user = {
    'name': 'testUser',
    'password': 'testPassword',
    'login': 'testLogin'
}

test_header = {
    'X-LOGIN': 'testLogin',
    'X-TOKEN': None
}


def scope_url(scope):
    '''Create url from scope.

    :param scope: scope name
    '''
    return '%s/%s' % (base_url, scope)


def create_user(user_profile=None):
    '''Create test user.'''
    user_profile = user_profile or test_user
    rv = requests.post(scope_url('user'), data=user_profile)
    if not rv.ok:
        raise Exception('Create user failed.')
    return rv.json()


def login_user():
    '''Login test user.'''
    rv = requests.post(scope_url('user/login'), data=test_user)
    if not rv.ok:
        raise Exception('Login user failed.')
    return rv.json()['token']


def create_tweet():
    '''Create a test tweet.'''
    test_tweet = {
        'text': 'hello world'
    }
    rv = requests.post(scope_url('t'), headers=test_header, data=test_tweet)
    if not rv.ok:
        raise Exception('Create tweet failed.')
    return rv.json()['Id']


class UserTest(unittest.TestCase):

    def testCreateUser(self):
        test_user = {
            'name': 'testUser',
            'password': 'testPassword',
            'login': 'testUser'
        }
        rv = requests.post(scope_url('user'), data=test_user)
        self.assertEqual(201, rv.status_code)

        rv = requests.post(scope_url('user'), data=test_user)
        self.assertEqual(409, rv.status_code)
        self.assertIn('error', rv.json())

        test_user['name'] = ''
        rv = requests.post(scope_url('user'), data=test_user)
        self.assertEqual(403, rv.status_code)

        test_user['password'] = '123'
        rv = requests.post(scope_url('user'), data=test_user)
        self.assertEqual(403, rv.status_code)

    def testReadUserProfile(self):
        rv = requests.get(scope_url('user/1'), headers=test_header)
        j = rv.json()
        self.assertEqual(200, rv.status_code)
        self.assertEqual(1, j['Id'])
        self.assertIn('t', j.keys())
        self.assertIsInstance(j['t'], list)
        #self.assertIn('Comments', j)

        rv = requests.get(scope_url('user/1'), headers={})
        self.assertEqual(401, rv.status_code)

        rv = requests.get(scope_url('user/123'), headers=test_header)
        self.assertEqual(404, rv.status_code)

    def testUpdateUserProfile(self):
        rv = requests.put(scope_url('user'), data=dict(name='test'),
                          headers=test_header)
        self.assertEqual(204, rv.status_code)
        rv = requests.get(scope_url('user/1'), headers=test_header)
        self.assertEqual(200, rv.status_code)
        self.assertEqual('test', rv.json()['Name'])

    def testLogout(self):
        rv = requests.get(scope_url('user/logout'), headers=test_header)
        self.assertEqual(200, rv.status_code)
        rv = requests.get(scope_url('user/1'), headers=test_header)
        self.assertEqual(401, rv.status_code)

        # Login again for later test.
        test_header['X-TOKEN'] = login_user()

    def testSelfProfile(self):
        rv = requests.get(scope_url('user/me'), headers=test_header)
        self.assertEqual(200, rv.status_code)
        j = rv.json()
        self.assertEqual(test_user['name'], j['Name'])
        self.assertNotIn('Password', j.keys())
        self.assertIn('t', j.keys())

        rv = requests.get(scope_url('user/me'))
        self.assertEqual(401, rv.status_code)

    def testUploadAvatar(self):
        rv = requests.post(scope_url('user/avatar'), headers=test_header)
        self.assertEqual(403, rv.status_code)

        rv = requests.post(scope_url('user/avatar'), headers=test_header,
                           files={'not-avatar': open(__file__, 'rb')})
        self.assertEqual(403, rv.status_code)

        rv = requests.post(scope_url('user/avatar'), headers=test_header,
                           files={'avatar': open(__file__, 'rb')})
        self.assertEqual(204, rv.status_code)

    def testUpdatePassword(self):
        origin = test_user['password']
        test_user['password'] = 'newPassword'
        payload = dict(oldPassword=origin, newPassword=test_user['password'])

        rv = requests.put(scope_url('user/password'), headers=test_header,
                          data=payload)
        self.assertEqual(204, rv.status_code)

        # Relogin user.
        test_header['X-TOKEN'] = login_user()

        rv = requests.put(scope_url('user/password'), headers=test_header,
                          data=payload)
        self.assertEqual(403, rv.status_code)

        payload['oldPassword'] = payload['newPassword']
        payload['newPassword'] = ''
        rv = requests.put(scope_url('user/password'), headers=test_header,
                          data=payload)
        self.assertEqual(403, rv.status_code)

        # Restore the password
        payload['newPassword'] = origin
        test_user['password'] = origin
        rv = requests.put(scope_url('user/password'), headers=test_header,
                          data=payload)
        self.assertEqual(204, rv.status_code)
        test_header['X-TOKEN'] = login_user()

    def testFriend(self):
        new_user = create_user({
            'name': 'new',
            'password': 'newpassword',
            'login': 'new'
        })

        # Add a friend.
        rv = requests.put(scope_url('user/friend/%s' % ('new')),
                          headers=test_header)
        self.assertEqual(204, rv.status_code)

        rv = requests.get(scope_url('user/me'), headers=test_header)
        r = rv.json()
        self.assertIsInstance(r['Friends'], list)
        self.assertIn(new_user['Id'], [i['Id'] for i in r['Friends']])

        rv = requests.get(scope_url('user/%d' % (new_user['Id'])),
                          headers=test_header)
        r = rv.json()
        self.assertIsInstance(r['Friends'], list)
        self.assertIn(1, [i['Id'] for i in r['Friends']])

        # Remove a friend.
        rv = requests.delete(scope_url('user/friend/%s' % ('new')),
                             headers=test_header)
        self.assertEqual(204, rv.status_code)

        rv = requests.get(scope_url('user/me'), headers=test_header)
        r = rv.json()
        self.assertIsInstance(r['Friends'], list)
        self.assertNotIn(new_user['Id'], [i['Id'] for i in r['Friends']])

        rv = requests.get(scope_url('user/%d' % (new_user['Id'])),
                          headers=test_header)
        r = rv.json()
        self.assertIsInstance(r['Friends'], list)
        self.assertNotIn(1, [i['Id'] for i in r['Friends']])


class TweetTest(unittest.TestCase):

    def testGetTimeline(self):
        rv = requests.get(scope_url('user/me'), headers=test_header)
        friendIds = [i['Id'] for i in rv.json()['Friends']]
        friendIds.append(1)

        create_tweet()

        rv = requests.get(scope_url('t'), headers=test_header)
        self.assertEqual(200, rv.status_code)
        j = rv.json()
        self.assertIsInstance(j, dict)
        self.assertIsInstance(j['t'], list)
        self.assertGreater(len(j['t']), 0)

        # Should only contains friends & user's tweets.
        for t in j['t']:
            self.assertIn(t['User']['Id'], friendIds)

        rv = requests.get(scope_url('t'))
        self.assertEqual(401, rv.status_code)

    def testCreateTweet(self):
        test_tweet = {
            'text': 'hello world'
        }
        rv = requests.post(scope_url('t'), headers=test_header,
                           data=test_tweet)
        self.assertEqual(201, rv.status_code)
        j = rv.json()
        id = j['Id']
        self.assertEqual(1, j['User']['Id'])
        rv = requests.get(scope_url('t/%d' % (id)), headers=test_header)
        self.assertEqual(200, rv.status_code)
        self.assertEqual(j['Text'], rv.json()['Text'])
        self.assertIsInstance(j['Comments'], list)
        self.assertIsInstance(j['Likes'], list)

        test_tweet['text'] = ''
        rv = requests.post(scope_url('t'), headers=test_header,
                           data=test_tweet)
        self.assertEqual(403, rv.status_code)

        rv = requests.post(scope_url('t'), headers=test_header,
                           files=dict(image=open(__file__, 'rb')))
        self.assertEqual(201, rv.status_code)

        rv = requests.post(scope_url('t'), headers=test_header,
                           data=test_tweet,
                           files=dict(image=open(__file__, 'rb')))
        self.assertEqual(201, rv.status_code)

        rv = requests.post(scope_url('t'), headers=test_header)
        self.assertEqual(403, rv.status_code)

    def testDeleteTweet(self):
        id = create_tweet()
        dest = scope_url('t/%d' % (id))
        rv = requests.delete(dest, headers=test_header)
        self.assertEqual(204, rv.status_code)
        rv = requests.get(dest, headers=test_header)
        self.assertEqual(404, rv.status_code)

    def testLikeTweet(self):
        id = create_tweet()
        dest = scope_url('t/%d/like' % (id))

        rv = requests.put(dest, headers=test_header)
        self.assertEqual(204, rv.status_code)
        rv = requests.get(scope_url('t/%d' % (id)), headers=test_header)
        self.assertIsNotNone(rv.json()['Likes'])

        rv = requests.put(dest, headers=test_header)
        self.assertEqual(204, rv.status_code)
        rv = requests.get(scope_url('t/%d' % (id)), headers=test_header)

    def testCreateComment(self):
        test_comment = {
            'content': 'hello world'
        }
        id = create_tweet()
        rv = requests.post(scope_url('t/%d/comment' % (id)),
                           headers=test_header, data=test_comment)
        self.assertEqual(201, rv.status_code)

        rv = requests.get(scope_url('t/%d' % (id)), headers=test_header)
        j = rv.json()
        self.assertEqual(1, len(j['Comments']))
        self.assertEqual(1, j['Comments'][0]['User']['Id'])
        self.assertEqual(test_comment['content'], j['Comments'][0]['Content'])


if __name__ == '__main__':
    create_user()
    test_header['X-TOKEN'] = login_user()
    unittest.main()
