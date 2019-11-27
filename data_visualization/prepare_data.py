import json
import sys

data_location = 'data_visualization/assets/data/data.json'
output_location = 'data_visualization/assets/data/data_subject_cat_group2.js'


def parseArgv():
    global data_location
    global output_location
    ln_arg = len(sys.argv)
    for i in range(1, ln_arg):
        s = sys.argv[i]
        if str.startswith(s, '-'):
            if s == '-i' and i+1 < ln_arg and not str.startswith(sys.argv[i+1], '-'):
                data_location = sys.argv[i+1]
            if s == '-o' and i+1 < ln_arg and not str.startswith(sys.argv[i+1], '-'):
                output_location = sys.argv[i+1]
    print(f'input: {data_location}; output: {output_location}')


def getData():
    with open(data_location, 'r', encoding="utf8") as f:
        lines = f.readlines()
    result = []
    for course in lines:
        c_json = json.loads(course)
        x = {}
        x['id'] = c_json['_id']
        x['provider'] = c_json['provider']

        x['categories'] = c_json['categories']
        x['subject'] = c_json['subject']
        x['schools'] = c_json['schools']
        x['teachers'] = c_json['teachers']
        try:
            x['rating'] = int(c_json['rating']['$numberInt'])
        except KeyError:
            x['rating'] = None
        try:
            x['review_count'] = int(c_json['rating']['$numberInt'])
        except KeyError:
            x['review_count'] = None
        try:
            x['language'] = c_json['details']['language']
        except KeyError:
            x['language'] = None

        result.append(x)
    return result


def saveJsFile(data):
    json_str = json.dumps(data, indent=2)
    json_str = 'modelDataAvailable(JSON.parse(`'+json_str+"""`
    ),{label: 'our data',file:'prepareddata_subject_cat_group_data2.jsonp'})
    """
    with open(output_location, 'w', encoding='utf8') as f:
        f.write(json_str)


def transformData(data):
    json_res = {"groups": []}
    subjects = {}

    for x in data:
        subject = x['subject']
        if subject not in subjects:
            subjects[subject] = set()

        for cat in x['categories']:
            subjects[subject].add(cat)
    id = 0
    for key in subjects:
        cats = subjects[key]
        group = {'label': key, 'weight': len(cats)*2, 'group': [], 'id': id}
        id += 1
        for cat in cats:
            c = {'label': cat, 'id': id, 'weight': 2}
            group['group'].append(c)
            id += 1
        json_res['groups'].append(group)

    return json_res


def main():
    parseArgv()

    res = getData()
    res = transformData(res)
    saveJsFile(res)

    print('done')


if __name__ == "__main__":
    main()