# How to Contribute

1. A new branch must be tagged with `WIP` until it is ready for CR.
   If the code doesn't pass review or has a merge conflict then the PR
   should be retagged until the code is ready to go back into CR.
2. Each commit must provide a short but specific description of the changes.
   Something like `Updated the range construct to properly scope index` or
   `Added unit-tests for the package creation` but not `unit-tests` or
   `Fixed the index problem`.
3. All public entities must be commented. If a parameter has any specific
   range, function, may not be null/nil etc. please add details about it
   in the comment. All private entities should also be commented but
   are not required to be commented.
4. All code should be auto-formatted.
5. The code may contain `TODO`s but should have the follow up branches
   defined. This may include follow up branches for commenting,
   adding examples or unit-tests, or fixes to know shortcuts or issues.   
6. Unit-tests should be added for all new features.
   Not all functions need to be specifically unit-tested.
7. The PR needs to have the branch in it's title.
   And the body should be similar to the following:
   ```
    [Issue #](path to issue, if there is one)

    ### Feature, Bug Fix, or Improvement
    A summary of the feature, steps to reproduce the bug,
    or a description of the improvement.

    ### Implementation
    - A list of changes made.

    ### Tech Debt
    - A list of anything not accomplished (should be labelled
      with TODOs) and anything that still needs to be done.
   ```
8. The PR must get at least one code review before being merged into master.
   The code review must check commenting, style, and consistency with
   designs. Check that feature(s) are implement and any tech debt has been
   called out or addressed. At least one reviewer must also code to confirm
   it builds and runs, and check that all unit-tests pass.
