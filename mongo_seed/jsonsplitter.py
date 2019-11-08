import json

inputfiles = ['courses.json']

for inputfile in inputfiles:
    with open('courses/data/'+inputfile) as json_file:
        data = json.load(json_file)   
       
        data = data["courses"]

        for key in data.keys():           
            category = data[key]
            
            for key in category.keys():
                newjson = category[key]
                newjson['_id'] = newjson['id']
                del newjson['id']

                with open('courses/data/output/'+str(key)+'.json', 'w') as out_file:          
                    json.dump(newjson,out_file, indent = 4)
