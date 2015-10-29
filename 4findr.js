//
// 2015
// Brennan Hoeting
// Oxford, OH
//

/**
 * Collections
 */

var Courses = new Mongo.Collection('courses');
var Classes = new Mongo.Collection('classes');
var Subjects = new Mongo.Collection('subjects');
var Professors = new Mongo.Collection('professors');

/**
 * |-----Courses---------|
 * | _id:        $string |
 * | gpa:        $number |
 * | title:      $string |
 * | short:      $string |
 * | number:     $number |
 * | subject_id: $number |
 * |---------------------|

 * |-----Classes------------|
 * | _id:           $string |
 * | gpa:           $number |
 * | course_id:     $number |
 * | professor_id:  $number |
 * |------------------------|
 *
 * |-----Subject------------|
 * | _id:           $string |
 * | gpa:           $number |
 * | short:         $string |
 * | number:        $number |
 * |------------------------|
 *
 * |---Professor---|
 * | _id:  $string |
 * | name: $number |
 * |---------------|
 */

/**
 * Config
 */

const DEPT_ID_URL = 'http://grdist.miamioh.edu/php/getDID.php?dept=%s';
const CLASSES_URL = 'http://grdist.miamioh.edu/php/getClasses.php?dept=%s&num=&inst=&from=2015&to=2015&iid=-1&did=%s&sem=&loc=O';
const DEPTARTMENTS = [
  { short: 'AAA', title: 'Asian/Asian American Studies'},
  { short: 'ACC', title: 'Accountancy'},
  { short: 'ACE', title: 'American Culture & English Prg'},
  { short: 'AER', title: 'Aeronautics'},
  { short: 'AES', title: 'Aerospace Studies'},
  { short: 'AMS', title: 'American Studies'},
  { short: 'ARB', title: 'Arabic'},
  { short: 'ARC', title: 'Architecture & Interior Design'},
  { short: 'ART', title: 'Art'},
  { short: 'ATH', title: 'Anthropology'},
  { short: 'BIO', title: 'Biology'},
  { short: 'BIS', title: 'Integrative Studies'},
  { short: 'BLS', title: 'Business Legal Studies'},
  { short: 'BOT', title: 'Botany'},
  { short: 'BTE', title: 'Business Technology'},
  { short: 'BUS', title: 'Business Analysis'},
  { short: 'BWS', title: 'Black World Studies'},
  { short: 'CAS', title: 'College of Arts and Science'},
  { short: 'CCA', title: 'College of Creative Arts'},
  { short: 'CEC', title: 'Col of Engineering & Computing'},
  { short: 'CHI', title: 'Chinese'},
  { short: 'CHM', title: 'Chemistry & Biochemistry'},
  { short: 'CIT', title: 'Comp Information Technology'},
  { short: 'CJS', title: 'Criminal Justice Studies'},
  { short: 'CLS', title: 'Classics'},
  { short: 'CMS', title: 'Comparative Media Studies'},
  { short: 'COM', title: 'Communication'},
  { short: 'CPB', title: 'Chem, Paper & Biomed Engineer'},
  { short: 'CPE', title: 'Chemical & Paper Engineering'},
  { short: 'CPS', title: 'Prof Studies & Appl Sciences'},
  { short: 'CRD', title: 'Civic and Regional Development'},
  { short: 'CSE', title: 'Comp Sci &Software Engineering'},
  { short: 'DST', title: 'Disability Studies'},
  { short: 'EAS', title: 'Engineering & App Science'},
  { short: 'ECE', title: 'Electrical & Computer Engineer'},
  { short: 'ECO', title: 'Economics'},
  { short: 'EDL', title: 'Educational Leadership'},
  { short: 'EDP', title: 'Educational Psychology'},
  { short: 'EDT', title: 'Teacher Education'},
  { short: 'EGM', title: 'Engineering Management'},
  { short: 'EHS', title: 'Education, Health and Society'},
  { short: 'ENG', title: 'English'},
  { short: 'ENT', title: 'Engineering Technology'},
  { short: 'ENV', title: 'Environmental Science'},
  { short: 'ESP', title: 'Entrepreneurship'},
  { short: 'FAS', title: 'Fashion Design'},
  { short: 'FIN', title: 'Finance'},
  { short: 'FRE', title: 'French'},
  { short: 'FST', title: 'Film Studies'},
  { short: 'FSW', title: 'Family Studies and Social Work'},
  { short: 'GEO', title: 'Geography'},
  { short: 'GER', title: 'German'},
  { short: 'GHS', title: 'Global Health Studies'},
  { short: 'GLG', title: 'Geology'},
  { short: 'GRK', title: 'Greek Language and Literature'},
  { short: 'GSC', title: 'Graduate School Community'},
  { short: 'GTY', title: 'Gerontology'},
  { short: 'HBW', title: 'Hebrew'},
  { short: 'HIN', title: 'Hindi'},
  { short: 'HON', title: 'Honors'},
  { short: 'HST', title: 'History'},
  { short: 'IDS', title: 'Interdisciplinary'},
  { short: 'IES', title: 'Environmental Sciences'},
  { short: 'IMS', title: 'Interactive Media Studies'},
  { short: 'ISA', title: 'Information Systems& Analytics'},
  { short: 'ITL', title: 'Italian'},
  { short: 'ITS', title: 'International Studies'},
  { short: 'JPN', title: 'Japanese'},
  { short: 'JRN', title: 'Journalism'},
  { short: 'KNH', title: 'Kinesiology and Health'},
  { short: 'KOR', title: 'Korean'},
  { short: 'LAS', title: 'Latin American Studies'},
  { short: 'LAT', title: 'Latin Language & Literature'},
  { short: 'LR:', title: 'Inst Learning in Retirement'},
  { short: 'LST', title: 'Liberal Studies'},
  { short: 'LUX', title: 'Luxembourg'},
  { short: 'MAC', title: 'Media and Culture'},
  { short: 'MBI', title: 'Microbiology'},
  { short: 'MGT', title: 'Management'},
  { short: 'MKT', title: 'Marketing'},
  { short: 'MME', title: 'Mechan & Manufact Engineering'},
  { short: 'MSC', title: 'Military Science'},
  { short: 'MTH', title: 'Mathematics'},
  { short: 'MUS', title: 'Music'},
  { short: 'NSC', title: 'Naval Science'},
  { short: 'NSG', title: 'Nursing'},
  { short: 'PCE', title: 'Paper & Chemical Engineering'},
  { short: 'PHL', title: 'Philosophy'},
  { short: 'PHY', title: 'Physics'},
  { short: 'PLW', title: 'Pre-Law Studies'},
  { short: 'PMD', title: 'Premedical Studies'},
  { short: 'POL', title: 'Political Science'},
  { short: 'POR', title: 'Portuguese'},
  { short: 'PSY', title: 'Psychology'},
  { short: 'REL', title: 'Religion, Comparative'},
  { short: 'RUS', title: 'Russian'},
  { short: 'SCA', title: 'School of Creative Arts'},
  { short: 'SJS', title: 'Social Justice Studies'},
  { short: 'SOC', title: 'Sociology'},
  { short: 'SPA', title: 'Speech Pathology & Audiology'},
  { short: 'SPN', title: 'Spanish'},
  { short: 'STA', title: 'Statistics'},
  { short: 'STC', title: 'Strategic Communication'},
  { short: 'THE', title: 'Theatre'},
  { short: 'UNV', title: 'University'},
  { short: 'WGS', title: 'Women, Gender & Sexuality Studies'},
  { short: 'WST', title: 'Western Program'},
  { short: 'ZOO', title: 'Zoology'}
];

/**
 * Server
 */

if (Meteor.isServer) {

  /**
   * Get the raw JSON data with the classes from the original source
   * @return {array} The raw data in JSON format
   */

  var getRawData = function () {
    let classes = [];
    for (dept of DEPTARTMENTS) {
      // Fetch the department ID for the next request
      let deptId = HTTP.get(sprintf(DEPT_ID_URL, dept.short)).data[0];
      if (!deptId) return;
      deptId = deptId.did;
      if (deptId == -1) return;

      // Fetch the classes for a certain deptartment
      let classesOfDept = HTTP.get(sprintf(CLASSES_URL, dept.short, deptId));
      if (!classesOfDept) return;
      classes = classes.concat(classesOfDept);
    }
    return classes;
  }

  Meteor.startup(function () {
    var data = getRawData();
    console.log(data);
  });
}

if (Meteor.isClient) {
  // counter starts at 0
  Session.setDefault('counter', 0);

  Template.hello.helpers({
    counter: function () {
      return Session.get('counter');
    }
  });

  Template.hello.events({
    'click button': function () {
      // increment the counter when button is clicked
      Session.set('counter', Session.get('counter') + 1);
    }
  });


}
