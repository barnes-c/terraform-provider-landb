schema_version = 1

project {
  license        = "GPL-3.0-or-later"
  copyright_year = 2025
  copyright_holder = "Christopher Barnes <christopher.barnes@cern.ch>"

  header_ignore = [
    # examples used within documentation (prose)
    "examples/**",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yml",

    # encrypted test fixtures
    "test/fixtures/**",
  ]
}