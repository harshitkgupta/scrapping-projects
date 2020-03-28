# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html
import pymongo
import configparser
class QuotePipeline(object):
    def __init__(self):
        config = configparser.ConfigParser()
        config.read('scrapy.cfg')
        host = config['mongodb']['host']
        port = int(config['mongodb']['port'])
        db = config['mongodb']['db']
        self.conn = pymongo.MongoClient(host, port)
        db_ref = self.conn[db]
        self.collection = db_ref['quotes']
    def process_item(self, item, spider):
        self.collection.insert(dict(item))
        return item
