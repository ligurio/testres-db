# testres-db

is a tool and library to import test results to a single database. It allows to
analyze your efficiency of testing, improve test execution and get better
vizualization of testing in your project. Queries can be executed with a
standard `sqlite` client or using [Jupyter Notebook](https://jupyter.org/).
Below there are some examples of SQL queries:

```sql
~$ sqlite testres-db.sqlite
SQLite version 2.8.17
Enter ".help" for instructions
sqlite> SELECT * FROM testrun TR JOIN test T ON TR.test_id = T.test_id JOIN status ST ON TR.status_id = ST.status_id WHERE ST.name = 'PASSED';
sqlite> SELECT * FROM testrun TR JOIN test T ON TR.test_id = T.test_id JOIN suite S ON T.suite_id = S.suite_id JOIN report R ON TR.report_id = R.report_id JOIN status ST ON TR.status_id = ST.status_id WHERE ST.name = 'PASSED';
```

## Building

```
$ go get ./...
$ go test -v ./...
$ cd cmd/testres-db && go build .
```

## How-To Use

First of all you should make sure test reports is available after every build
on continuous integration. For example GitLab CI [allows](https://docs.gitlab.com/ee/ci/junit_test_reports.html) to store JUnit reports as artifacts when Travis CI not.

Secondly you should describe configuration file in YAML format which contains
information about your project.

```yaml
projects:
- name: criu
  backends:
   - type: jenkins
     branch: master
     pipeline: ""
     base: https://ci.openvz.org/
     username: anonymous
     secret: anonymous
```

By default `testres-db` uses `testres-db.yaml` as configuration file and
`testres.sqlite` as an SQLite database file. You can change default values
using command-line options.

## Contribution

There are numerous awesome [Continuous Integration
services](https://github.com/ligurio/awesome-ci) which are not integrated with
`testres-db`. Feel free to make a patches and bring support for them.

## Authors

Developed with passion by [Sergey Bronnikov](https://bronevichok.ru/) and great
open source [contributors](https://github.com/ligurio/testres-db/contributors).
