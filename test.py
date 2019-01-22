#!/usr/local/bin/python3.6

from py2neo import Node, Relationship
 
a = Node('Person', name='Alice')
b = Node('Person', name='Bob')
r = Relationship(a, 'KNOWS', b)
print(a, b, r)