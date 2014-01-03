#!/usr/bin/env python

# Python library contains codepages compiled from Unicode.org. This script uses this tables for work indirectly :)

import os, sys
from cStringIO import StringIO
import encodings

def load(nm):
    fnm = os.path.join(os.path.dirname(__file__), nm)
    f = open(fnm, "rt")
    t = f.read()
    f.close()
    return t

# Convert all codepages files to one Go module:
TABLES = {
     "ascii": [],
     "cp037": [],
     "cp1006": [],
     "cp1026": [],
     "cp1140": [],
     "cp1250": [],
     "cp1251": [],
     "cp1252": [],
     "cp1253": [],
     "cp1254": [],
     "cp1255": [],
     "cp1256": [],
     "cp1257": [],
     "cp1258": [],
     "cp424": [],
     "cp437": [],
     "cp500": [],
     "cp720": [],
     "cp737": [],
     "cp775": [],
     "cp850": [],
     "cp852": [],
     "cp855": [],
     "cp856": [],
     "cp857": [],
     "cp858": [],
     "cp860": [],
     "cp861": [],
     "cp862": [],
     "cp863": [],
     "cp864": [],
     "cp865": [],
     "cp866": [],
     "cp869": [],
     "cp874": [],
     "cp875": [],
     "cp932": [],
     "cp949": [],
     "cp950": [],
     "euc_jis_2004": [],
     "euc_jisx0213": [],
     "euc_jp": [],
     "euc_kr": [],
     "gb18030": [],
     "gb2312": [],
     "gbk": [],
     "hp_roman8": [],
     "iso2022_jp": [],
     "iso2022_jp_1": [],
     "iso2022_jp_2": [],
     "iso2022_jp_2004": [],
     "iso2022_jp_3": [],
     "iso2022_jp_ext": [],
     "iso2022_kr": [],
     "iso8859_1": [],
     "iso8859_10": [],
     "iso8859_11": [],
     "iso8859_13": [],
     "iso8859_14": [],
     "iso8859_15": [],
     "iso8859_16": [],
     "iso8859_2": [],
     "iso8859_3": [],
     "iso8859_4": [],
     "iso8859_5": [],
     "iso8859_6": [],
     "iso8859_7": [],
     "iso8859_8": [],
     "iso8859_9": [],
     "koi8_r": [],
     "koi8_u": [],
     "latin_1": [],
     "mac_arabic": [],
     "mac_centeuro": [],
     "mac_croatian": [],
     "mac_cyrillic": [],
     "mac_farsi": [],
     "mac_greek": [],
     "mac_iceland": [],
     "mac_latin2": [],
     "mac_roman": [],
     "mac_romanian": [],
     "mac_turkish": []
}

def gen_test(nm):
    us = []
    os = []
    for i in range(32, 256):
        c = chr(i)
        try:
            us.append(c.decode(nm))
            os.append(c)
        except:
            pass
    s = u''.join(us)
    orig = "".join(os)

    for cp in [ "utf-8", "utf-16le", "utf-16be", "utf-32le", "utf-32be" ]:
        fnm = "test_%s_%s.txt" % (nm, cp)
        f = open(fnm, "wb")
        f.write(s.encode(cp))
        f.close()
    f = open("test_%s_ORIG.txt" % nm, "wb")
    f.write(orig)
    f.close()

# Process aliases:
aliases = encodings.aliases.aliases
for a in aliases:
    if a not in TABLES:
        if aliases[a] in TABLES:
            TABLES[aliases[a]].append(a)

def pytable(enc):
    " This function creates unicode string with all encoded characters "
    r = { }
    for i in range(256):
        if i < 128:
            r[chr(i)] = i
        else:
            c = chr(i)
            try:
                code = ord(c.decode(enc))
            except:
                code = 0
            r[chr(i)] = code
    return r

tblnum = 0
def table2string(tbl):
    global tblnum
    tblnum += 1
    S = StringIO()
    S.write("var tbl_%d = [256]rune{\n\t0x0000" % tblnum)
    for i in range(1, 256):
        S.write(',')
        if i % 16 == 0:
            S.write('\n\t')
        if chr(i) in tbl:
            S.write("0x%04x" % tbl[chr(i)])
        else:
            S.write("0x0000")
    S.write("}\n")
    return 'tbl_%d' % tblnum, S.getvalue()

def reversetable(tbl):
    global tblnum
    tblnum += 1
    S = StringIO()

    l = []
    for i in range(256):
        l.append( (chr(i), tbl[chr(i)]) )
    l.sort(lambda a, b: a[1] < b[1])

    S.write("var tbl_%d = [256]pair{" % tblnum)
    for i in range(256):
        if i > 0:
            S.write(',')
        if i % 4 == 0:
            S.write('\n\t')
        else:
            S.write(' ')
        S.write("{ 0x%02x, 0x%04X }" % (i, tbl[chr(i)]))
    S.write("}\n")
    return 'tbl_%d' % tblnum, S.getvalue()

def gen_tables():
    f = open(os.path.join(os.path.dirname(__file__), "tables8.go"), "wt")
    pre = load("pre.go")
    f.write(pre)

    post = load("post.go")
    tbls = {}
    for t in TABLES:
        tbl = pytable(t)
        nm1, a = table2string(tbl)
        nm2, b = reversetable(tbl)
        f.write(a)
        f.write('\n')
        f.write(b)
        f.write('\n')
        tbls[t] = (nm1, nm2)
    # Make table with names:
    f.write("var names = [...]tbls{")
    comma = ""
    for t in TABLES:
        f.write('%s\n\ttbls{"%s", %s, %s}' % (comma, t, tbls[t][0], tbls[t][1]))
        comma = ","
        for n in TABLES[t]:
            f.write('%s\n\ttbls{"%s", %s, %s}' % (comma, n, tbls[t][0], tbls[t][1]))
    f.write("}\n")

    f.write(post)
    f.close()
    print "DONE"

def gen_tests():
    print "Generate tests for several characters encodings..."
    TEST = []
    TEST.append("""#!/bin/sh

T=/tmp/test-goconv.$$
GOCONV=../test/test

error()
{
    echo "Error: $*" 1>&2
    exit 1
}

test()
{
    "$GOCONV" -f "$1" -t "$2" -o "$T" "test_$1_ORIG.txt" || error "$1 -> $2"
    diff test_$1_$2.txt "$T" || error "$1 -> $2: invalid"
    rm -f "$T"
}

""")

    for t in TABLES:
        gen_test(t)
        for cp in [ "utf-8", "utf-16le", "utf-16be", "utf-32le", "utf-32be" ]:
            TEST.append("test %s %s\n" % (t, cp))

    f = open("test.sh", "wt")
    f.write("".join(TEST))
    f.close()

if __name__ == '__main__':
    if len(sys.argv) > 1 and sys.argv[1] == 'test':
        gen_tests()
    else:
        gen_tables()


