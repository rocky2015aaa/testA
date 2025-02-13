
# HashOCR Server
### Prerequisite
docker
### How To Run
- Start server

At current repo, run `make setup`
- Test

After running the server as described:
1. Import `HashOCR.postman_collection.json` from `internal/testdata` into Postman.
2. Switch the test file in your directory to the appropriate form.

Or

Use curl command
```
1) curl --location 'http://localhost:9081/api/v1/pdf2json?invoice_type=tax' --form 'file=@"{file_location}/27350AA.pdf"'
2) curl --location 'http://localhost:9081/api/v1/pdf2json?invoice_type=bakerty' --form 'file=@"{file_location}/bakerty.pdf"'
3) curl --location 'http://localhost:9081/api/v1/pdf2json?invoice_type=winners' --form 'file=@"{file_location}/WinnersInvoice.pdf"'
4) curl --location 'http://localhost:9081/api/v1/pdf2json?invoice_type=tambaram' --form 'file=@"{file_location}/TambaramInvoice.pdf"'
5) curl --location 'http://localhost:9081/api/v1/pdf2json?invoice_type=rp' --form 'file=@"{file_location}/RPInvoice.pdf"'
```
- Remove Docker Container and Image

`make stop`

`make clean`

`make clean-image`

Or

`make clean-all`