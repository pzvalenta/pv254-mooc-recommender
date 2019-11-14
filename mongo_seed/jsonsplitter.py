import json

# takes a string in format 123 or 123.4k and converts it to int
def myStringToInt(string):
    fres = 0.0
    numericpart = ""
    unit = ""
    for c in string:
        if c.isdigit() or c == ".":          
            numericpart += c
        else:
            unit += c
    
    if ( len(numericpart) + len(unit) ) != len(string):
        raise SyntaxError
    
    fres = float(numericpart)

    if unit == "":
        pass
    elif unit == "k":
        fres *= 1000
    else:
        raise SyntaxError
    
    return int(fres) 

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
                newjson['interested_count'] = myStringToInt(newjson['interested_count'])
                newjson['review_count'] = myStringToInt(newjson['review_count'])
                
                # convert any string into array containing that string
                if type(newjson['details']['start date']) is str:
                    newjson['details']['start_date'] = [newjson['details']['start date']]
                else:
                    newjson['details']['start_date'] = newjson['details']['start date']
                del newjson['details']['start date']
                
                with open('courses/data/output/'+str(courseID)+'.json', 'w') as out_file:          
                    json.dump(newjson,out_file, indent = 4)
