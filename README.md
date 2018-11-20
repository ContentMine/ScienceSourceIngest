Tool for importing openXML format papers into ScienceSource
---------

Run ./ingest -help to see the command options.





How to connect the took to the wiki
================

Make a consumer on wikibase
------------------

To make a consumer you will need an account that has an email address If you're using the wikibase/wikibase docker image and haven't set up email then use the maintenance script resetUserEmail.php and don't forget --no-reset-password

Now log in as the admin user

Now navigate to http://<your wikibase>/Special:OAuthConsumerRegistration

Select request a token for a new consumer

Register a consumer by filling in the form:

Application name (e.g. 'Quickstatements)
Version (you can leave this as 1.0)
Description (e.g. 'Quickstatements')
callback URL: http://{your quickstatements host}/api.php
Check "Allow consumer to specify a callback in requests"
Request permission for:
High-volume editing
Edit existing pages
Create, edit, and move pages
Check agree and click "propose consumer"
Make a note of the details on the following page
