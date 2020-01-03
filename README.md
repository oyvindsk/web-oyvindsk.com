# web-oyvindsk.com

## About
Simple go blog that serves http://oyvindsk.com and its alias http://odots.org/

I'm curentltly using it to test Google Cloud Run. It used to run on App Engine Standard, and I might convert it back in the future.
Both are container based PaaS solutions. Run is newer, is Docker based and offers more flexibility. It does lack some App Engine's features though (such as traffic splitting).
Both a re fully managed (yay!) and can scale up "infinitely" and down to zero when there's no traffic.

This blog "engine" reads posts and pages on startup, executes the templates and stores the resulting html in memory.
The blogposts are also go templates, to have a little more freedome. Might be a bad idea.

The static files are served directly from Google Cloud Storage (GCS).


## Todo
Too much to list =/

## Lisence
Go code (*.go files) are MIT / BSD 2-clause. Files in the static directory have their own licenses, see the files themselves. The css files in static/ have an unknown license, but was pretty much copied from http://www.alexedwards.net/.
