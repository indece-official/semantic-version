#!/bin/bash -e

WORKDIR=../build/test-integration/tmp/
PROGRAM_DEFAULT=../dist/bin/semantic-version
PROGRAM=$(realpath "${PROGRAM:-$PROGRAM_DEFAULT}")

before() {
    rm -rf "$WORKDIR"
    mkdir -p "$WORKDIR"
    mkdir -p "$WORKDIR/repo"
    cd $WORKDIR/repo
}

assertVersion() {
    VERSION=$($PROGRAM get-version)
    if [[ $? -ne 0 ]] ; then
        echo "ERROR: $PROGRAM returned error code $?"

        exit 1
    fi

    if [[ "$VERSION" != "$1" ]] ; then
        echo "ERROR: Expected version $1, got $VERSION"

        exit 1
    fi
}

assertChangelogLines() {
    CHANGELOG=$($PROGRAM get-changelog)
    if [[ $? -ne 0 ]] ; then
        echo "ERROR: $PROGRAM returned error code $?"

        exit 1
    fi

    COUNT_RELEVANT_LINES=$(echo "$CHANGELOG" | grep '^* ' | wc -l)
    if [[ "$COUNT_RELEVANT_LINES" != "$1" ]] ; then
        echo $CHANGELOG
        echo "ERROR: Expected $1 changes in changelog, got $COUNT_RELEVANT_LINES"

        exit 1
    fi
}

testSimple() {
    echo "Testing simple repository"

    git init > /dev/null
    
    echo "1" > "testfile.txt"
    git add . > /dev/null
    git commit -m "Initial commit" > /dev/null
    assertVersion "v1.0.0"
    assertChangelogLines 0

    echo "2" > "testfile.txt"
    assertVersion "v1.0.0"
    assertChangelogLines 0

    git add . > /dev/null
    git commit -m "feat: Some change" > /dev/null
    assertVersion "v1.0.0"
    assertChangelogLines 1

    git tag v1.0.0
    assertVersion "v1.0.0"
    assertChangelogLines 0

    echo "3" > "testfile.txt"
    git add . > /dev/null
    git commit -m "feat: Some change 2" > /dev/null
    assertVersion "v1.1.0"
    assertChangelogLines 1

    git tag -a v1.1.0 -m "Version 1.1.0"
    assertVersion "v1.1.0"
    assertChangelogLines 0

    echo "4" > "testfile.txt"
    git add . > /dev/null
    git commit -m "fix: Some fix" > /dev/null
    assertVersion "v1.1.1"
    assertChangelogLines 1

    echo "Success"
}

testBranches() {
    echo "Testing repository with branches"

    git init > /dev/null
    
    echo "1" > "testfile1.txt"
    git add . > /dev/null
    git commit -m "Initial commit" > /dev/null
    assertVersion "v1.0.0"
    assertChangelogLines 0

    git checkout -b feat/test1
    assertVersion "v1.0.0-feat_test1.0"
    assertChangelogLines 0

    echo "2" > "testfile2.txt"
    git add . > /dev/null
    git commit -m "feat: Some change feat 1" > /dev/null
    assertVersion "v1.0.0-feat_test1.0"
    assertChangelogLines 1
    
    echo "3" > "testfile3.txt"
    git add . > /dev/null
    git commit -m "fix: Some change feat 1" > /dev/null
    assertVersion "v1.0.0-feat_test1.0"
    assertChangelogLines 2

    git checkout master
    assertVersion "v1.0.0"
    assertChangelogLines 0

    git tag -a v1.1.0 -m "Version 1.1.0"
    assertVersion "v1.1.0"
    assertChangelogLines 0

    echo "4" > "testfile4.txt"
    git add . > /dev/null
    git commit -m "feat: Some change master 1" > /dev/null
    assertVersion "v1.2.0"
    assertChangelogLines 1

    git tag -a v1.2.0 -m "Version 1.2.0"
    assertVersion "v1.2.0"
    assertChangelogLines 0

    git checkout feat/test1
    assertVersion "v1.2.0-feat_test1.0"
    assertChangelogLines 2

    git tag v1.2.0-feat_test1.0
    assertVersion "v1.2.0-feat_test1.1"
    assertChangelogLines 2

    echo "5" > "testfile5.txt"
    git add . > /dev/null
    git commit -m "fix: Some change feat 1" > /dev/null
    assertVersion "v1.2.0-feat_test1.1"
    assertChangelogLines 3

    git checkout master
    assertVersion "v1.2.0"
    assertChangelogLines 0

    git merge feat/test1 -m "Merge feat/test1"
    assertVersion "v1.3.0"
    assertChangelogLines 3

    echo "Success"
}

main() {
    before
    testSimple

    before
    testBranches
}

main