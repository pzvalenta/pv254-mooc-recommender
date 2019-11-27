import json
import sys

data_location = 'data_visualization/data.json'
output_location = 'data_visualization/prepared_data.json'

ln_arg = len(sys.argv)
for i in range(1, ln_arg):
    s = sys.argv[i]
    if str.startswith(s, '-'):
        if s == '-i' and i+1 < ln_arg and not str.startswith(sys.argv[i+1], '-'):
            data_location = sys.argv[i+1]
        if s == '-o' and i+1 < ln_arg and not str.startswith(sys.argv[i+1], '-'):
            output_location = sys.argv[i+1]


print(f'input: {data_location}; output: {output_location}')

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


with open(output_location, 'w', encoding='utf8') as f:
    json.dump(result, fp=f, indent=2)

print('done')