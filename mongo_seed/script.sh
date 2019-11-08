for f in courses/data/output/*.json ; do mongoimport --host=mongo --collection=courses --db=mydb --file="$f"; done
