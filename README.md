Solr Command Line Interface (solrcli) is a command line program that emulates
the query feature of the Solr Admin web page that Solr provides out of the box.

This is a work in progress in *very* early stages.

To build
```
git clone https://github.com/hectorcorrea/solrcli.git
go get
go build
```

To use:
```
$ solrcli http://localhost:8983/solr/your-core
```

