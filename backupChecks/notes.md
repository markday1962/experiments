
```
FOR a IN assignee
  FILTER SUM(FOR v IN 1..1 OUTBOUND a GRAPH "groupingGraph" RETURN 1) == 0
  RETURN a
```

otherwise there's an audit log, so could check the latest timestamp on that maybe

but a better way might be to look at the individual collections and check the max ID of each of those


https://html.developreference.com/article/10391849/Query+code+for+finding+highest+value+of+collections%3F

FOR c IN organisation
LIMIT 1
SORT c._key DESC
RETURN c

[
  {
    "_key": "77961",
    "_id": "organisation/77961",
    "_rev": "_aX8wa-6--_",
    "name": "White & Wyckoff Mfg Company",
    "isNpe": false,
    "size": 14,
    "sizeActive": 0,
    "numOffensiveDisputes": 0,
    "numDefensiveDisputes": 0,
    "address": "Afghanistan"
  }
]

FOR c IN assignee
LIMIT 1
SORT c._key DESC
RETURN c

[
  {
    "_key": "1308",
    "_id": "assignee/1308",
    "_rev": "_aX8Xdra--J",
    "address": "FR",
    "defensiveDisputeIds": [],
    "isNpe": false,
    "name": "COMMISSARIAT ENERGIE ATOMIQUE",
    "offensiveDisputeIds": [],
    "patFamIds": [
      2605,
      79580
    ]
  }
]
