with open('data.json', 'r', encoding="utf8") as f:
    lines = f.readlines()
with open('data_fixed.json', 'w', encoding="utf8") as f:
    f.write('[\n')
    for i, line in enumerate(lines):
        if i < len(lines)-1:
            line = line.rstrip('\n') + ',\n'
        f.write(line)
    f.write(']\n')
