import argparse
import sys

CHANGELOG_FILE = "CHANGELOG.md"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("run", type=str, help="the job to run")
    parser.add_argument("--tag", type=str, help="the Git tag to work with")
    args = parser.parse_args()

    if args.run == "check-changelog":
        check_changelog(args.tag)
    elif args.run == "print-changelog":
        print_changelog(args.tag)

    sys.exit(0)


def check_changelog(git_tag):
    """
    Check if a new release tag is mentioned in the changelog.

    For a release tag like `v1.2.3`, the changelog has to contain a
    release section called `[1.2.3]`. If the release isn't mentioned
    in the changelog, exit with an error.
    """
    # Cut off the `v` prefix to get the actual release number.
    search_expr = "[{0}]".format(git_tag[1:])

    with open(CHANGELOG_FILE) as changelog:
        content = changelog.read()
        if search_expr not in content:
            msg = """You're trying to create a new release tag {0}, but that release is not mentioned
in the changelog. Add a section called {1} to {2} and try again.""" \
                .format(git_tag, search_expr, CHANGELOG_FILE)

            sys.exit(msg)


def print_changelog(git_tag):
    """
    Print the changelog for the given release tag by reading the
    changelog file. If the release tag does not exist as a release
    number in the changelog, the output will be empty.
    """
    start = "## [{0}]".format(git_tag[1:])
    # The ## [Unreleased] heading will be ignored.
    unreleased = "## [Unreleased]"
    end = "## ["

    capturing = False
    output = ""

    with open(CHANGELOG_FILE) as changelog:
        lines = changelog.readlines()
        for line in lines:
            # Start capturing if the line contains our release heading.
            if start in line and unreleased not in line:
                capturing = True
                continue
            # Stop capturing if we've reached the end, i.e. the next heading.
            if end in line and capturing:
                break
            if capturing:
                output += line

    print(output)


main()
