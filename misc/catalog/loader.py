import os
import sys
import json
import pymongo

def main(host='localhost', db=None):
    dir_path = os.path.dirname(os.path.realpath(__file__))

    if db is None:
        db = pymongo.Connection(host=host).keyfu
    catalog = db.catalog

    data = None
    with open(os.path.join(dir_path, 'data.json')) as f:
        data = json.load(f)

    insert_data = []
    for k, v in data.items():
        v['value'] = k
        v['sort'] = v['name'].lower()
        catalog.update({'value': k}, v, upsert=True)

if __name__ == '__main__':
    if len(sys.argv) > 1:
        main(sys.argv[1])
    else:
        main()
