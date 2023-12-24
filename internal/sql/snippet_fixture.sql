INSERT INTO snippet (title, content, created, expired) VALUES (
'An old silent pond',
'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
localtimestamp,
(localtimestamp + INTERVAL '365 DAYS')
);

INSERT INTO snippet (title, content, created, expired) VALUES (
'The Lazy FOX',
'The quick brown fox jump over the lazy badger',
localtimestamp,
(localtimestamp + INTERVAL '7 DAYS')
);