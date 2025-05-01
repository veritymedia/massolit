import type { TeacherRaw, ExamRaw } from "../objects";

const teachers: TeacherRaw[] = [
  {
    name: "Simon",
    subjects: ["English"],
    availabilities: [
      {
        dow: "3",
        start: "9:00",
        end: "9:55",
      },
      {
        dow: "3",
        start: "11:10",
        end: "12:05",
      },
    ],
  },
  {
    name: "Jhon",
    subjects: [""],
    availabilities: [
      {
        dow: "1",
        start: "9:00",
        end: "9:55",
      },
      {
        dow: "1",
        start: "11:10",
        end: "12:05",
      },
    ],
  },
];
