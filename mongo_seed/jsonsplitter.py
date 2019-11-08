import json

inputfiles = ['art-and-design_courses.json',
'business_courses.json',
'cs_courses.json',
'data-science_courses.json',
'education_courses.json',
'engineering_courses.json',
'health_courses.json',
'humanities_courses.json',
'maths_courses.json',
'personal-development_courses.json',
'programming-and-software-development_courses.json',
'science_courses.json',
'social-sciences_courses.json']

for inputfile in inputfiles:
    with open('courses/data/'+inputfile) as json_file:
        data = json.load(json_file)
        for key in data.keys():
            newjson = data[key]
            del newjson['id']
            newjson["name"] = str(key)

            with open('output/'+str(key)+'.json', 'w') as out_file:          
                json.dump(newjson,out_file, indent = 4)
