Tool for importing openXML format papers into a wikibase instance
============


This tool takes a JSON list of papers that are both on WikiData and PubMedCentral, processes them into HTML and does simple dictionary based text mining, the output of which is pushed to a wikibase instance.

This tool requires you have xsltproc installed to let it process the openXML papers.

Run ./ingest -help to see the command options and details. There are six things you need to provide typically (explained in more detail below):

* -feed [file path] - this is a JSON file that contains a list of the papers as fetched from WikiData
* -output [directory path] - this is a directory where the tool will store its working state
* -urlbase [http(s)://wikibase.server.name] - This should be the protocol and hostname of your Wikibase server
* -oauth [file path] - a JSON file containing the Consumer and Access information for your Wikibase
* -dictionaries [directory path] - this is a directory where the dictionaries of words to be annotated are found
* -xsltproc [file path] - this is the location of the xsltproc tool. Defaults to "/usr/bin/xsltproc"


Paper Feed
--------

There is an example feed of papers in the repository, please see that as a full example. It was generated from Wikidata using a SPARQL query of the form:

```
SELECT DISTINCT ?item ?itemLabel ?pmcid ?journalLabel ?title ?date ?licenseLabel ?mainsubjectLabel
  WHERE
        {
          ?item wdt:P31 wd:Q13442814 .
          ?item wdt:P5008 wd:Q55439927 .
          ?item wdt:P932 ?pmcid .
          ?item wdt:P1433 wd:Q3359737.
          ?item wdt:P1433 ?journal .
          ?item wdt:P1476 ?title .
          ?item wdt:P577 ?date .
          ?item wdt:P275 ?license .
          ?item wdt:P921 ?mainsubject .
          ?mainsubject wdt:P361* wd:Q18123741 .

        SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],en". }

        }
 LIMIT 100
```

This should be exported from wikidata as a JSON feed.


Output
------

The output directory is where ScienceSourceIngest will store its state, and it is recommend you use the same output directory for multiple runs to the same wikibase target server, but a different output directory per wikibase target server.


URL base
--------

This should be the prefix of the URL for the wikibase instance you're uploading to. E.g., to upload to Science Source you'd set it to https://sciencesource.wmflabs.org, or for a local test instance to http://localhost:8181 and so forth.


OAuth information
------------------

For ScienceSourceIngest to talk to your wikibase target server you need to authenticate with it. On the target wikibase server you should navigate to /wiki/Special:OAuthConsumerRegistration/propose and pick the following options:

* Application name - set to something you'll remember for this use
* Consumer version - doesn't matter, so set it to v1.0 or such
* Application description - set to something you'll remember
* This consumer is for use only by [USERNAME] - set this on

The rest can remain at defaults. That last one is the most important, so please ensure you select that. If not then wikibase will require you to authorise the client via a web interface which will not work.

For Applicable grants select the following:

* Basic rights
* High-volume editing
* Edit existing pages
* Edit protected pages
* Create, edit, and move pages
* Protect and unprotect pages

When you click done you will find a page with the following 4 strings on it:

* Consumer Token
* Consumer Secret
* Access Token
* Access Secret

If you don't see these you probably forgot to tick "This consumer is for use only by [USERNAME]". You should take these values and put them in a JSON file like so:

```
{
    "consumer": {
        "key": "633d4025c53c4179ba7260a801a6aee3",
        "secret": "9b4403abf1e9e7fbb208866081df0e5f5770d322"
    },
    "access": {
        "token": "9f5b9e2e20b655299c699e5172dd9ba1",
        "secret": "3c31dd374458e2ce34a4b5163f3ae231cb304d3c"
    }
}
```

You then pass this file as a parameter when you start ScienceSourceIngest.

Dictionaries
------------

The annotations that ScienceSourceIngest finds in the papers are based on the dictionaries supplied here. There are sample dictionaries in the project dictionaries folder.


Usage notes
-----------

Please note that uploading data in bulk can be slow - annotations require a lot of items to be created and properties to be set in the Wikibase instance, and each call will take around a second to complete on a remote server, which means papers can take a minute or so to upload fully.

If you re-run the program with the same input feed and output directory then it should safely resume upload from where it left off and not re-upload anything it had already uploaded.


License
============

This software is copyright Content Mine Ltd 2018, and released under the Apache 2.0 License.


Dependencies
============

Relies on https://github.com/ContentMine/wikibase
