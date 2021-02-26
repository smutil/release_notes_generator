# Release Metrics Generator ![example workflow](https://github.com/smutil/release_notes_generator/actions/workflows/build-actions.yml/badge.svg) ![example workflow](https://github.com/smutil/release_notes_generator/actions/workflows/release-actions.yml/badge.svg)

CLI to generate release notes based on tag and git commit log in html and json format with below details.
1. Release Name
2. Change Volume
3. Change Leadtime
4. Author
5. Release Date
6. Commit list with leadtime


Usage
-----
 step 1. download release_notes_generator from <a href=https://github.com/smutil/release_notes_generator/releases>releases</a>. 
 
 step 3. execute the release_notes_generator as shown below. --releaseVersion is optional. latest release will be used if --releaseVersion is not provided.
 
 ```
 ./release_notes_generator --application findevops --gitRepo  https://github.com/xxxx/xxxxxxxx.git -gitCred username:password --releaseVersion v1.1
 ```
 step 4. ReleaseNotes.html and ReleaseNotes.json will be generated in same location.

 ![Alt text](docs/images/release_notes.png?raw=true "Title")

