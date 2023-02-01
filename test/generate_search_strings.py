import itertools
import random

people_csv = "db/people.csv"
search_data_csv = "test/search_data.csv"


def substr(string, min_len=2):
    if len(string) <= min_len:
        yield string
        return

    for l in range(min_len, len(string) + 1):
        yield string[:l]
        # for shift in range(0, len(string) - l + 1):
        #     yield string[shift:shift + l].lower()


search_data = list()
with open(people_csv) as f:
    for line in f:
        name, _, _ = line.split(',')
        last_name, first_name = name.split(' ')

        def iterate():
            return itertools.product(substr(last_name), substr(first_name))

        for last, first in iterate():
            search_data.append(','.join((last, first)))


with open(search_data_csv, 'w') as f:
    random.shuffle(search_data)
    for line in search_data[:10000]:
        f.write(line)
        f.write('\n')
