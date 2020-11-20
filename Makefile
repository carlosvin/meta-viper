
md:
	asciidoctor -b docbook README.adoc
	iconv -t utf-8 README.xml | pandoc -f docbook -t markdown_mmd --highlight-style=pygments --wrap=none | iconv -f utf-8 > README.md

unittest:
	go test -coverprofile=coverage.out -v -coverpkg ./... ./...
	go tool cover -html=coverage.out