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
    VERSION=$($PROGRAM $2 $3 $4 $5 get-version)
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
    CHANGELOG=$($PROGRAM $2 $3 $4 $5 get-changelog)
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

    LAST_GIT_HASH=$(git rev-parse --short HEAD)
    git checkout $LAST_GIT_HASH
    echo "2a" > "testfile.txt"
    assertVersion "UNKNOWN"
    assertChangelogLines 0

    assertVersion "v1.0.0" -git-branch master
    assertChangelogLines 1 -git-branch master

    git checkout "testfile.txt"
    git checkout master
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
    assertVersion "v1.2.0-feat_test1.1" -debug
    assertChangelogLines 2 -debug

    echo "5" > "testfile5.txt"
    git add . > /dev/null
    git commit -m "fix: Some change feat 1" > /dev/null
    assertVersion "v1.2.0-feat_test1.1" -debug
    assertChangelogLines 3 -debug

    git checkout master
    assertVersion "v1.2.0"
    assertChangelogLines 0

    git merge feat/test1 -m "Merge feat/test1"
    assertVersion "v1.3.0" -debug
    assertChangelogLines 3 -debug

    echo "Success"
}

testDevelopReleaseSimple() {
    echo "Testing simple develop/release repository"

    git init > /dev/null
    
    echo "1" > "testfile.txt"
    git add . > /dev/null
    git commit -m "Initial commit" > /dev/null
    assertVersion "v1.0.0"
    assertChangelogLines 0

    git tag v1.0.0

    git checkout -b develop

    echo "1\n\n2" > "testfile.txt"
    git add . > /dev/null
    git commit -m "feat: Change 1" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0

    git checkout master
    assertVersion "v1.0.0"
    assertChangelogLines 0

    git merge develop --no-ff --commit -m "Merge develop in master"
    assertVersion "v1.1.0"
    assertChangelogLines 1

    git tag v1.1.0
    
    git checkout develop

    echo "1\n\n2\n\n3" > "testfile.txt"
    git add . > /dev/null
    git commit -m "fix: Change 2" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0

    git checkout master
    assertVersion "v1.1.0"
    assertChangelogLines 0

    git merge develop --no-ff --commit -m "Merge develop in master"
    assertVersion "v1.1.1"
    assertChangelogLines 1

    echo "Success"
}

testDevelopReleaseComplex() {
    echo "Testing complex develop/release repository"

    cat >./semanticversion.yaml <<EOL
branches:
  - branch_pattern: master
    release_channel: FINAL
    version_pattern: v{major}.{minor}.{patch}

  - branch_pattern: release.*
    release_channel: FINAL
    version_pattern: v{major}.{minor}.{patch}

  - branch_pattern: testing.*
    release_channel: BETA
    version_pattern: v{major}.{minor}.{patch}-beta.{build}
EOL

    git init > /dev/null

    echo "1" > "testfile1.txt"
    git add . > /dev/null
    git commit -m "Initial commit" > /dev/null
    assertVersion "v1.0.0"
    assertChangelogLines 0

    git tag v1.0.0

    git checkout -b release-1.x
    echo "2" > "testfile2.txt"
    git add . > /dev/null
    git commit -m "fix: 2" > /dev/null
    assertVersion "v1.0.1"
    assertChangelogLines 1

    git tag v1.0.1

    git checkout -b testing-1.x
    echo "3" > "testfile3.txt"
    git add . > /dev/null
    git commit -m "fix: 3" > /dev/null
    assertVersion "v1.0.2-beta.0"
    assertChangelogLines 1

    git checkout -b develop-1.x

    echo "4" > "testfile4.txt"
    git add . > /dev/null
    git commit -m "feat: 4" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0
    
    echo "5" > "testfile5.txt"
    git add . > /dev/null
    git commit -m "fix: 5" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0

    git checkout testing-1.x

    git merge develop-1.x --no-ff --commit -m "Merge develop-1.x in testing-1.x"
    assertVersion "v1.1.0-beta.0"
    assertChangelogLines 3

    git tag v1.1.0-beta.0

    git checkout release-1.x

    git checkout -b release-2.x

    git checkout -b testing-2.x

    git checkout -b develop-2.x

    echo "6" > "testfile6.txt"
    git add . > /dev/null
    git commit -m "break: 6" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0
    
    echo "7" > "testfile7.txt"
    git add . > /dev/null
    git commit -m "fix: 7" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0

    git checkout testing-2.x

    git merge develop-2.x --no-ff --commit -m "Merge develop-2.x in testing-2.x"
    assertVersion "v2.0.0-beta.0"
    assertChangelogLines 2

    git tag v2.0.0-beta.0
    
    git checkout release-2.x

    git merge testing-2.x --no-ff --commit -m "Merge testing-2.x in release-2.x"

    assertVersion "v2.0.0"
    assertChangelogLines 2

    git tag v2.0.0
    
    git checkout develop-2.x

    echo "8" > "testfile8.txt"
    git add . > /dev/null
    git commit -m "feat: 8" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0
    
    echo "9" > "testfile9.txt"
    git add . > /dev/null
    git commit -m "fix: 9" > /dev/null
    assertVersion "UNKNOWN"
    assertChangelogLines 0

    git checkout testing-2.x

    git merge develop-2.x --no-ff --commit -m "Merge develop-2.x in testing-2.x"
     
    assertVersion "v2.1.0-beta.0"
    assertChangelogLines 2

    git tag v2.1.0-beta.0
    
    git checkout release-2.x

    git merge testing-2.x --no-ff --commit -m "Merge testing-2.x in release-2.x"

    assertVersion "v2.1.0"
    assertChangelogLines 2

    git tag v2.1.0

    git checkout release-1.x
    assertVersion "v1.0.1"
    assertChangelogLines 0
    
    git checkout testing-1.x
    assertVersion "v1.1.0-beta.1"  -debug
    assertChangelogLines 0

    git checkout release-1.x

    git merge testing-1.x --no-ff --commit -m "Merge testing-1.x in release-1.x"
    assertVersion "v1.1.0"
    assertChangelogLines 3

    git tag v1.1.0
    assertVersion "v1.1.0"
    assertChangelogLines 0

    echo "Success"
}

main() {
    before
    testSimple

    before
    testBranches

    before
    testDevelopReleaseSimple

    before
    testDevelopReleaseComplex
}

main