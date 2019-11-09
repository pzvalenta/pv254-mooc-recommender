import json

inputfiles = ['courses.json']

for inputfile in inputfiles:
    with open('courses/data/'+inputfile) as json_file:
        data = json.load(json_file)   
       
        data = data["courses"]

        for category in data.keys():                     
            courses = data[category]
            
            for courseID in courses.keys():
                newjson = courses[courseID]
                newjson['_id'] = newjson['id']
                del newjson['id']

                with open('courses/data/output/'+str(courseID)+'.json', 'w') as out_file:          
                    json.dump(newjson,out_file, indent = 4)
