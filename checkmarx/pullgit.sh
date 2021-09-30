#!/bin/bash
rm -Rf gftn
mkdir gftn
cd gftn
git clone git@github.ibm.com:gftn/gftn-services.git
cd ..
java -Xmx1024m -jar $PWD/CxConsolePlugin-CLI-8.80.0-20180806-1131.jar Scan -v -ProjectName "\CxServer\SP\Company\World Wire\WorldWire" -CxServer "http://52.206.144.150" -cxUser "srinivas.veerasamy@ibm.com" -cxPassword "Ch3ckm@rx$" -locationtype folder -locationpath $PWD/gftn
