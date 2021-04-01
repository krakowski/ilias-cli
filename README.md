This is a small console application with the primary goal of managing and
grading exercise tasks within the ILIAS eLearning platform.

## :wrench: &nbsp; Installation

Precompiled binaries are available for Linux, macOS and Windows from
the [GitHub Releases page](https://github.com/krakowski/ilias-cli/releases).

> :satellite: &nbsp; **Update Checks**
>
> The console application automatically checks at startup if a newer
> version is available and notifies the user if it is present.

## :notebook: &nbsp; Workflow

The console application is intended for tutors as well as for the supervising teaching staff
and thus offers functions for distributing as well as correcting submissions.

### Teaching staff

For each exercise sheet to be corrected, a correction bundle must be compiled for the tutors,
with which they can download their assigned submissions by means of the console application.
The basic procedure can be described as follows.

1. Create a file named `.tutors.yml` in which the weekly correction time in hours of all participating tutors is entered.
   ```
   - id: tutor100
     hours: 9
  
   - id: tutor101
     hours: 7
  
   - id: tutor102
     hours: 5
   ```
   The field 'id' reflects the university ID belonging to the respective tutor. This file can be reused.


2. Open the corresponding ILIAS exercise object in the browser and note down the URL's `ref_id` parameter. Since this
   value will not change for the lifetime of the object, this step can be skipped in the future.


3. Open the "Submissions and Grades" view and extract the value of the `<select>` option of the desired exercise unit
   using the Dev Tools of the browser. The `id` of the `<select>` DOM element in the HTML source code is `ass_id`.


4. Invoke the console tool with the arguments `exercises distribute <ref_id> <ass_id>` where the parameters `ref_id`
   and `ass_id` should be replaced with the previously determined values. This will create a file named `.workspace.yml`
   which contains the correction assignments.


5. Create a correction template named `CORRECTION.tmpl` which will be provided to the tutors. An example of the content
   may look like the following.
   ```
   student: {{.Student}}
   points: 0
   corrected: false
   correction: |
  
     <b>!! Rückfragen bitte unter Angabe der Korrekturnummer '{{.Student}}' an {{.Tutor}}@hhu.de !!</b>
  
     <b>Aufgabe 1.1 (0/7)</b>
    
     -
  
     <b>Aufgabe 1.2 (0/18)</b>
    
     -
  
     <b>Aufgabe 1.3 (0/18)</b>
    
     -
  
     <b>Aufgabe 1.4 (0/7)</b>
    
     -
   ```

   The template variables `{{.Student}}` as well as `{{.Tutor}}` are replaced by the console tool during the correction.

6. Make the files `.workspace.yml` and `CORRECTION.tmpl` available to the tutors, for example, bundled in an archive.

### Tutors

1. Unpack or place the files `.workspace.yml` and `CORRECTION.tmpl` in a common folder and call the console tool inside
   this folder with the arguments `workspace init`. This creates a separate subfolder for each correction, which contains
   the submission as well as the populated correction template.

2. Correct the submissions and then set the `points` and `corrected` fields. If the `corrected` field is not set to
   `true` for a submission, the corrections cannot be uploaded. Make sure to use a dot (`.`) as separator for decimal numbers.

3. Use the console tool's `workspace status` command to check if there are any outstanding corrections.

4. Invoke the console tool with the arguments `workspace upload` to upload all corrections. Alternatively, it is also
   possible to upload individual corrections by using the `--only` flag with comma separated IDs of the corrections to be uploaded.


## :triangular_ruler: &nbsp; Configuration Options

To avoid having to enter the login data every time the console tool is invoked, they can alternatively be stored
within environment variables for the duration of a session.

> :warning: &nbsp; **Credentials**
>
> It is not a good idea to store the login data inside these variables permanently, because they will be available in plain text.

|      Name     |   Type   | Description                                    |
|:-------------:|:--------:|------------------------------------------------|
| `ILIAS_USER`  | `string` |  The username used for logging into ILIAS      |
| `ILIAS_PASS`  | `string` |  The specified user's login password           |

`ILIAS_USER` and `ILIAS_PASS` are not required, since they are interactively queried from the console when not available.

## :scroll: &nbsp; License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.