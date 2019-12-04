import json
import sys

data_location = 'data_visualization/assets/data/data.json'
output_location = 'data_visualization/assets/data/'


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
    print(f'input: {data_location}; output-path: {output_location}')


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


def saveJsFile(data, name):
    print(f'saving {name} to {output_location+name}.js')
    json_str = json.dumps(data, indent=2)
    json_str = 'modelDataAvailable('+json_str+"""
    ,{label: 'our data',file:'"""+name+""".js}.jsonp'})
    """
    with open(output_location+name+'.js', 'w', encoding='utf8') as f:
        f.write(json_str)


def subject_cat_data(data):
    json_res = {"groups": []}
    subjects = {}
    subjects_course_count = {}
    cat_course_count = {}

    for x in data:
        subject = x['subject']
        if subject not in subjects:
            subjects[subject] = {}
            subjects_course_count[subject] = 0
        subjects_course_count[subject] += 1
        for cat in x['categories']:
            if cat not in subjects[subject]:
                subjects[subject][cat] = 0
            subjects[subject][cat] += 1
    id = 0
    for key in subjects:
        cats = subjects[key]
        group = {'label': key,
                 'weight': subjects_course_count[key], 'groups': [], 'id': id}
        id += 1
        for cat in cats:
            c = {'label': cat, 'id': id, 'weight': subjects[key][cat]}
            group['groups'].append(c)
            id += 1
        json_res['groups'].append(group)

    return json_res


def languages_data(data):
    languages = {}
    for x in data:
        lan = x['language']
        if lan not in languages:
            languages[lan] = 0
        languages[lan] += 1
    json_res = {"groups": []}
    id = 0
    for lan in languages:
        lan_c = languages[lan]
        group = {'label': f'{lan}: {lan_c}',
                 'weight': lan_c, 'id': id}
        id += 1
        json_res['groups'].append(group)

    return json_res


def provider_subject_data(data):
    json_res = {"groups": []}
    providers = {}
    provider_course_count = {}
    provider_subjects_count = {}

    for x in data:
        subject = x['subject']
        provider = x['provider']
        if provider not in providers:
            providers[provider] = set()
            provider_course_count[provider] = 0

        provider_course_count[provider] += 1
        if provider not in provider_subjects_count:
            provider_subjects_count[provider] = {}

        if subject not in provider_subjects_count[provider]:
            provider_subjects_count[provider][subject] = 0

        provider_subjects_count[provider][subject] += 1

        providers[provider].add(subject)

    id = 0

    for key in providers:
        subs = providers[key]
        group = {'label': key,
                 'weight': provider_course_count[key], 'groups': [], 'id': id}
        id += 1
        for sub in subs:
            c = {'label': sub, 'id': id,
                 'weight': provider_subjects_count[key][sub]}
            group['groups'].append(c)
            id += 1
        json_res['groups'].append(group)
    return json_res


def subject_provider_data(data):
    json_res = {"groups": []}
    subjects = {}
    subjects_course_count = {}
    subjects_providers_count = {}

    for x in data:
        subject = x['subject']
        provider = x['provider']

        if subject not in subjects:
            subjects[subject] = set()
            subjects_course_count[subject] = 0
        subjects_course_count[subject] += 1

        if subject not in subjects_providers_count:
            subjects_providers_count[subject] = {}

        if provider not in subjects_providers_count[subject]:
            subjects_providers_count[subject][provider] = 0

        subjects_providers_count[subject][provider] += 1

        subjects[subject].add(provider)
    id = 0

    for key in subjects:
        provs = subjects[key]
        group = {'label': key,
                 'weight': subjects_course_count[key], 'groups': [], 'id': id}
        id += 1
        for prov in provs:
            c = {'label': prov, 'id': id,
                 'weight': subjects_providers_count[key][prov]}
            group['groups'].append(c)
            id += 1
        json_res['groups'].append(group)

    return json_res


def main():
    parseArgv()

    data = getData()
    res = subject_cat_data(data)
    saveJsFile(res, "subject_cat_groups_data")
    # res2 = subject_provider_data(data)
    # saveJsFile(res2, "subject_provider_groups_data")
    # res3 = languages_data(data)
    # saveJsFile(res3, "languages_data")
    # res3 = provider_subject_data(data)
    # saveJsFile(res3, "provider_subject_groups_data")
    print('done')


if __name__ == "__main__":
    main()
