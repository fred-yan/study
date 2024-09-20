class Entry:
    def __init__(self):
        self.key = ''
        self.level = 0
        self.children = None

original = [
    "Nonmetals",
    "    Hydrogen",
    "    Carbon",
    "    Nitrogen",
    "    Oxygen",
    "Inner Transitionals",
    "    Lanthanides",
    "        Europium",
    "        Cerium",
    "    Actinides",
    "        Uranium",
    "        Plutonium",
    "        Curium",
    "Alkali Metals",
    "    Lithium",
    "    Sodium",
    "    Potassium",
]

def populate_entries(original_str):
    original_list = []
    prefix = '    '
    for item in original_str:
        level = 0
        while str(item).startswith(prefix):
            item = item[len(prefix):]
            level = level + 1
        else:
            entry = Entry()
            entry.key = str(item).strip()
            entry.level = level
            entry.children = []
            add_entry(original_list, entry)
    return original_list

def add_entry(list_of_entries: [], data: Entry):
    if data.level == 0:
        list_of_entries.append(data)
    else:
        level = 0
        append_entry = list_of_entries[-1]
        list_of_append_entry = []
        while data.level > level:
            list_of_append_entry = append_entry.children
            if list_of_append_entry:
                append_entry = list_of_append_entry[-1]
                level = append_entry.level
            else:
                break

        list_of_append_entry.append(data)


def visit_process_list(p_list):
    if p_list:
        sort_p_list = sorted(p_list, key=lambda x: x.key)
        for p in sort_p_list:
            print(p.level * '    ' + p.key)
            c_list = p.children
            visit_process_list(c_list)
    else:
        return


process_list = populate_entries(original)
visit_process_list(process_list)
