export interface ExamRaw {
  subject: string;
  start: string; // format hh:mm or ISO format 2025-05-06T08:00:00
  duration: string; // format HH:mm or mm (minutes)
  room: string; // where the exam will take place
  date?: string; // formatted dd/mm/yyy
  id?: string | number; // optional id field
  examCode?: string; // optional exam code field
}

interface Exam extends ExamRaw {
  id: string | number;
  startMinutes: number;
  endMinutes: number;
  dayOfWeek: number;
  proctors: {
    teacher: string;
    start: string;
    end: string;
  }[];
}

export interface TeacherRaw {
  name: string;
  schedule: string[]; // format ["0-1", "1-4"]
  subjects: string[];
}
interface Teacher {
  name: string;
  subjects: string[];
  schedule: { day: number; period: number }[];
  assignedMinutes: number;
}

interface ExamWithProctors extends ExamRaw {
  proctors: {
    teacher: string;
    start: string;
    end: string;
  }[];
}

interface SchoolTimetable {
  [key: number]: { start: string; end: string };
}

export interface Configuration {
  schoolTimetable: SchoolTimetable;
  parallel_teachers: number;
}

export function assignTeachersToExams(
  examsRaw: ExamRaw[],
  teachersRaw: TeacherRaw[],
  config: Configuration,
): ExamWithProctors[] {
  // ====== UTILITY FUNCTIONS ======

  function timeToMinutes(timeStr: string): number {
    // Handle ISO format
    if (timeStr.includes("T")) {
      const timePart = timeStr.split("T")[1];
      const [hours, minutes] = timePart.split(":").map(Number);
      return hours * 60 + minutes;
    }
    // Handle regular time format
    const [hours, minutes] = timeStr.split(":").map(Number);
    return hours * 60 + minutes;
  }

  function minutesToTime(minutes: number): string {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours.toString().padStart(2, "0")}:${mins.toString().padStart(2, "0")}`;
  }

  function parseDuration(durationStr: string): number {
    if (durationStr.includes(":")) {
      // Format is HH:mm
      const [hours, minutes] = durationStr.split(":").map(Number);
      return hours * 60 + minutes;
    } else {
      // Format is just minutes
      return parseInt(durationStr);
    }
  }

  function getDayOfWeekFromSlashDate(dateStr: string): number {
    const [day, month, year] = dateStr.split("/").map(Number);
    const date = new Date(year, month - 1, day);
    return date.getDay();
  }

  function getDayOfWeekFromISODate(dateStr: string): number {
    const date = new Date(dateStr);
    return date.getDay();
  }

  function getDateFromISODateTime(dateTimeStr: string): string {
    return dateTimeStr.split("T")[0];
  }

  function getLessonsSpanned(
    startMin: number,
    endMin: number,
    timetable: Configuration["schoolTimetable"],
  ): number[] {
    const lessons = [];
    for (const [periodStr, times] of Object.entries(timetable)) {
      const period = Number(periodStr);
      const periodStart = timeToMinutes(times.start);
      const periodEnd = timeToMinutes(times.end);

      if (startMin < periodEnd && endMin > periodStart) {
        lessons.push(period);
      }
    }
    return lessons;
  }

  function teachesSubject(
    teacherSubjects: string[],
    examSubject: string,
  ): boolean {
    const examSubjectLower = examSubject.toLowerCase();
    return teacherSubjects.some(
      (subject) =>
        examSubjectLower.includes(subject.toLowerCase()) ||
        subject.toLowerCase().includes(examSubjectLower),
    );
  }

  function isTeacherAvailable(
    teacher: Teacher,
    dayOfWeek: number,
    startMin: number,
    endMin: number,
  ): boolean {
    const lessons = getLessonsSpanned(startMin, endMin, config.schoolTimetable);
    return lessons.every(
      (lesson) =>
        !teacher.schedule.some(
          (slot) => slot.day === dayOfWeek && slot.period === lesson,
        ),
    );
  }

  function calculateLessonUtilization(
    startMin: number,
    endMin: number,
    dayOfWeek: number,
    teacher: Teacher,
  ): number {
    let usedTime = 0;
    let totalAvailableTime = 0;

    for (const slot of teacher.schedule) {
      if (slot.day !== dayOfWeek) continue;

      const lessonStart = timeToMinutes(
        config.schoolTimetable[slot.period].start,
      );
      const lessonEnd = timeToMinutes(config.schoolTimetable[slot.period].end);

      if (startMin < lessonEnd && endMin > lessonStart) {
        const overlapStart = Math.max(startMin, lessonStart);
        const overlapEnd = Math.min(endMin, lessonEnd);
        const overlapMinutes = overlapEnd - overlapStart;
        const lessonMinutes = lessonEnd - lessonStart;

        usedTime += overlapMinutes;
        totalAvailableTime += lessonMinutes;
      }
    }

    // Higher score = better utilization of whole lessons
    return totalAvailableTime > 0 ? usedTime / totalAvailableTime : 0;
  }

  // ====== MAIN ALGORITHM ======

  const parallel_teachers = config.parallel_teachers || 1;

  // Process teachers data
  const teachers: Teacher[] = teachersRaw.map((t) => ({
    name: t.name,
    subjects: t.subjects.map((s) => s.toLowerCase()),
    schedule: t.schedule.map((s) => {
      const [day, period] = s.split("-").map(Number);
      return { day, period };
    }),
    assignedMinutes: 0,
  }));

  // Process exams data
  const exams: Exam[] = examsRaw.map((exam, index) => {
    const startMinutes = timeToMinutes(exam.start);
    const durationMinutes = parseDuration(exam.duration);

    let dayOfWeek: number = 0;
    let examDate: string | undefined = exam.date;

    // Handle date formats
    if (exam.start.includes("T")) {
      // ISO format with date included in start time
      dayOfWeek = getDayOfWeekFromISODate(exam.start);
      examDate = getDateFromISODateTime(exam.start);
    } else if (exam.date && exam.date.includes("/")) {
      // Slash date format
      dayOfWeek = getDayOfWeekFromSlashDate(exam.date);
    } else if (exam.date && exam.date.includes("-")) {
      // ISO date format without time
      dayOfWeek = getDayOfWeekFromISODate(exam.date);
    }

    return {
      ...exam,
      id: exam.id !== undefined ? exam.id : index, // Keep original ID if exists
      startMinutes,
      endMinutes: startMinutes + durationMinutes,
      dayOfWeek,
      proctors: [],
      date: examDate,
    };
  });

  // Group exams by date and room
  const examsByDateAndRoom = new Map<string, Map<string, Exam[]>>();

  for (const exam of exams) {
    if (!exam.date) {
      console.warn(`Exam ${exam.id} has no date information`);
      continue;
    }

    if (!examsByDateAndRoom.has(exam.date)) {
      examsByDateAndRoom.set(exam.date, new Map());
    }

    const roomsMap = examsByDateAndRoom.get(exam.date)!;
    if (!roomsMap.has(exam.room)) {
      roomsMap.set(exam.room, []);
    }

    roomsMap.get(exam.room)!.push(exam);
  }

  // Process each date and room
  for (const [date, roomsMap] of examsByDateAndRoom.entries()) {
    for (const [room, roomExams] of roomsMap.entries()) {
      // Create timeline of all start and end points
      const timePoints = new Set<number>();
      for (const exam of roomExams) {
        timePoints.add(exam.startMinutes);
        timePoints.add(exam.endMinutes);
      }

      const sortedTimePoints = Array.from(timePoints).sort((a, b) => a - b);

      // Process each time slice
      for (let i = 0; i < sortedTimePoints.length - 1; i++) {
        const startTime = sortedTimePoints[i];
        const endTime = sortedTimePoints[i + 1];

        // Find active exams in this time slice
        const activeExams = roomExams.filter(
          (exam) =>
            exam.startMinutes <= startTime && exam.endMinutes >= endTime,
        );

        if (activeExams.length === 0) continue;

        // Determine day of week based on the date format
        let dayOfWeek: number;
        if (date.includes("/")) {
          dayOfWeek = getDayOfWeekFromSlashDate(date);
        } else {
          dayOfWeek = getDayOfWeekFromISODate(date);
        }

        const firstLessonStart = timeToMinutes(config.schoolTimetable[0].start);

        // Handle early morning exams
        if (startTime < firstLessonStart) {
          const earlyEnd = Math.min(firstLessonStart, endTime);

          // Assign DEFAULT TEACHER to early time period
          for (const exam of activeExams) {
            exam.proctors.push({
              teacher: "DEFAULT TEACHER",
              start: minutesToTime(startTime),
              end: minutesToTime(earlyEnd),
            });
          }

          // If time slice extends beyond school start time, continue with next iteration
          if (endTime <= firstLessonStart) continue;

          // Process the remainder of the time slice that's during school hours
          const remainingSliceStart = firstLessonStart;

          // Get subjects for this exam group
          const examSubjects = new Set(
            activeExams.map((e) => e.subject.toLowerCase()),
          );

          // Find eligible teachers for this time slice
          const eligibleTeachers = teachers.filter((teacher) => {
            // Check availability
            const available = isTeacherAvailable(
              teacher,
              dayOfWeek,
              remainingSliceStart,
              endTime,
            );
            if (!available) return false;

            // Check subject compatibility
            const canProctor = Array.from(examSubjects).every(
              (subject) => !teachesSubject(teacher.subjects, subject),
            );

            return canProctor;
          });

          // Score teachers based on lesson utilization and workload balance
          const teachersWithScore = eligibleTeachers.map((teacher) => ({
            teacher,
            score:
              calculateLessonUtilization(
                remainingSliceStart,
                endTime,
                dayOfWeek,
                teacher,
              ) -
              teacher.assignedMinutes / 1000, // Balance factor
          }));

          // Sort by score (higher is better)
          teachersWithScore.sort((a, b) => b.score - a.score);

          // Select teachers for assignment
          const neededTeachers = Math.min(
            parallel_teachers,
            activeExams.length,
          );
          const selectedTeachers = teachersWithScore
            .slice(0, neededTeachers)
            .map((t) => t.teacher);

          // Apply assignments
          for (const exam of activeExams) {
            for (const teacher of selectedTeachers) {
              exam.proctors.push({
                teacher: teacher.name,
                start: minutesToTime(remainingSliceStart),
                end: minutesToTime(endTime),
              });

              // Update teacher workload
              teacher.assignedMinutes += endTime - remainingSliceStart;
            }
          }
        } else {
          // Regular in-school time slice

          // Get subjects for this exam group
          const examSubjects = new Set(
            activeExams.map((e) => e.subject.toLowerCase()),
          );

          // Find eligible teachers
          const eligibleTeachers = teachers.filter((teacher) => {
            // Check availability
            const available = isTeacherAvailable(
              teacher,
              dayOfWeek,
              startTime,
              endTime,
            );
            if (!available) return false;

            // Check subject compatibility
            const canProctor = Array.from(examSubjects).every(
              (subject) => !teachesSubject(teacher.subjects, subject),
            );

            return canProctor;
          });

          // Score teachers based on lesson utilization and workload balance
          const teachersWithScore = eligibleTeachers.map((teacher) => ({
            teacher,
            score:
              calculateLessonUtilization(
                startTime,
                endTime,
                dayOfWeek,
                teacher,
              ) -
              teacher.assignedMinutes / 1000, // Balance factor
          }));

          // Sort by score (higher is better)
          teachersWithScore.sort((a, b) => b.score - a.score);

          // Select teachers for assignment
          const neededTeachers = Math.min(
            parallel_teachers,
            activeExams.length,
          );
          const selectedTeachers = teachersWithScore
            .slice(0, neededTeachers)
            .map((t) => t.teacher);

          if (selectedTeachers.length < neededTeachers) {
            console.warn(
              `Not enough eligible teachers for time slot ${minutesToTime(startTime)}-${minutesToTime(endTime)} on ${date} in room ${room}`,
            );
          }

          // Apply assignments
          for (const exam of activeExams) {
            for (const teacher of selectedTeachers) {
              exam.proctors.push({
                teacher: teacher.name,
                start: minutesToTime(startTime),
                end: minutesToTime(endTime),
              });

              // Update teacher workload
              teacher.assignedMinutes += endTime - startTime;
            }
          }
        }
      }
    }
  }

  // Convert to output format
  return exams.map(
    ({ startMinutes, endMinutes, dayOfWeek, ...rest }) => rest,
  ) as ExamWithProctors[];
}
