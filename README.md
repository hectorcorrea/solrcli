Solr Command Line Interface (`solrcli`) is a text based user interface to submit
queries to Solr directly from a Linux terminal.

Technically you can do anything that this tool does via cURL commands, but at
times crafting the right cURL command for a Solr query (with all the different
encoding options) can be trick. With `solrcli` you just enter the Solr
parameters as you would in the native "Solr Admin" web page and execute them.

This is a work in progress in *very* early stages.

## Compiling and running
```
git clone https://github.com/hectorcorrea/solrcli.git
go get
go build
./solrcli http://localhost:8983/solr/your-core
```

When running it more or less looks like this:
```
./solrcli http://localhost:8983/solr/your-core

==============================================================================
[h]elp | [q]uery | [f]acet | f[l] | e[x]ecute | [s]tart | r[o]ws | [Q]uit
==============================================================================
> x
2020/02/20 16:24:30 Solr HTTP GET: http://localhost:8983/solr/your-core/select?q=%2A&defType=edismax&debug=false&
{
	"responseHeader": {
		"status": 0,
		"QTime": 0,
		"params": {
			"q": "*",
			"defType": "edismax",
			"debug": "false"
		}
	},
	"response": {
		"numFound": 4296,
		"start": 0,
		"docs": [
            ...
        ]
    }
}
```

You can customize the query (`q`), facet field (`f`), list of fields to return
(`l`), and so on via the options in the menu bar.


## Executable

If you don't care about the source code, download the executable for your operating
system, and follow the instructions on Release tab.

