import {
  assignTeachersToExams,
  type Configuration,
  type ExamRaw,
  type TeacherRaw,
} from "./v2";
import { exams } from "./data/june2025";

const config: Configuration = {
  parallel_teachers: 1,
  schoolTimetable: {
    0: { start: "9:00", end: "9:54" },
    1: { start: "9:55", end: "10:49" },
    2: { start: "10:50", end: "11:09" },
    3: { start: "11:10", end: "12:04" },
    4: { start: "12:05", end: "12:59" },
    5: { start: "13:00", end: "13:39" },
    6: { start: "13:40", end: "14:34" },
    7: { start: "14:35", end: "15:30" },
  },
};

// const exams: ExamRaw[] = [
//   {
//     date: "08/05/2025",
//     duration: "55",
//     room: "room1",
//     start: "9:00",
//     subject: "english",
//   },
//   // { date: '08/05/2025', duration: '55', room: 'room1', start: '9:00', subject: 'science' },
//   // { date: '09/05/2025', duration: '60', room: 'room1', start: '9:00', subject: 'english' },
// ];

const teachers: TeacherRaw[] = [
  // { name: 'Teacher 1', schedule: ['0-0', '1-2', '1-1', '1-0'], subjects: ['biology'] },
  {
    schedule: [
      "0-0",
      "0-1",
      "0-3",
      "1-4",
      "1-5",
      "1-6",
      "2-6",
      "2-3",
      "2-1",
      "3-0",
      "3-3",
      "3-4",
      "3-6",
      "3-7",
    ],
    name: "Simon",
    subjects: ["english"],
  },
  { name: "Teacher 2", schedule: ["0-2", "0-1"], subjects: ["biology"] },
];

const res = assignTeachersToExams(exams, teachers, config);

console.log("RESULTS ", res);
