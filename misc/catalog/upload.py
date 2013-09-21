import libcloud.security
libcloud.security.VERIFY_SSL_CERT = True

import os
from glob import glob
from libcloud.storage.types import Provider
from libcloud.storage.providers import get_driver

connection = get_driver(Provider.CLOUDFILES_US)('silas', '717d0b3f2c9e6c7938a06922c557413d')

container = connection.get_container('KeyFu')

objects = {}

for obj in container.list_objects():
    objects[obj.name] = obj

names = [path[4:] for path in glob('img/*.png')]

for name in names:
    if name not in objects:
        print 'Uploading %s...' % name
        container.upload_object(os.path.join('img', name), name, extra={'content_type': 'image/png'})
    else:
        print 'Skipping %s...' % name
