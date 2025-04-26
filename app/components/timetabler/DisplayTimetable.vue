<template>
  <div class="rounded-lg text-foreground">
    <h1 class="mb-6 text-2xl font-bold">Exam Schedule</h1>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center py-10">
      <div
        class="w-10 h-10 border-t-2 border-b-2 border-indigo-500 rounded-full animate-spin"
      ></div>
    </div>

    <!-- Empty state -->
    <div v-else-if="localWeeklyExams.length === 0" class="py-10 text-center">
      No exams scheduled at this time.
    </div>

    <!-- Exams by week -->
    <div v-else class="flex flex-col w-full h-full">
      <div
        v-for="(week, weekIndex) in localWeeklyExams"
        :key="weekIndex"
        class="mb-8"
      >
        <div class="flex items-center mb-4">
          <div class="w-3 h-3 mr-2 bg-[purple] rounded-full"></div>
          <h2 class="text-xl font-semibold">
            {{ week.weekLabel }}
          </h2>
        </div>

        <!-- Week exams cards -->
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
          <div
            v-for="(examGroup, examIndex) in week.examGroups"
            :key="examIndex"
            class="p-4 transition-shadow duration-200 border-l-4 rounded-lg shadow-sm max-h-min bg-background bg-opacity-35 hover:shadow-md"
            :class="getBorderClass(examGroup)"
          >
            <!-- Subject(s) and completion badge -->
            <div class="flex items-start justify-between mb-2">
              <div>
                <h3 class="text-lg font-medium capitalize">
                  {{ getSubjectsLabel(examGroup) }}
                </h3>
                <span v-if="examGroup.exams.length > 1" class="text-xs">
                  {{ examGroup.exams.length }} combined exams
                </span>
              </div>

              <span
                v-if="examGroup.exams.length > 1"
                class="px-2 py-1 text-xs font-medium text-blue-800 bg-blue-100 rounded-full"
              >
                Multiple
              </span>
            </div>

            <!-- Time and room info -->
            <div class="flex flex-col gap-2">
              <div>
                <!-- <Icon name="material-symbols:delete" /> -->
                <span class="text-lg font-bold">{{
                  formatTime(examGroup.exams[0].bookedSegments[0].start)
                }}</span>
              </div>
              <div class="flex flex-col gap-2">
                <div>
                  <span class="pr-2 text-xl">📅</span>
                  {{ formatDate(examGroup.start) }} <br />
                </div>
                <div>
                  <span class="pr-2 text-xl">⏳</span>
                  {{ formatTime(getLatestEndTime(examGroup)) }} ({{
                    getDurationRange(examGroup)
                  }})
                </div>

                <div>
                  <span class="pr-2 text-xl">🏢</span>
                  {{ examGroup.room }}
                </div>
              </div>
            </div>

            <!-- Teachers -->
            <div v-if="getAllTeachers(examGroup).length > 0" class="mt-3">
              <h4 class="mb-1 text-sm font-medium">Proctored by:</h4>
              <div class="flex flex-wrap">
                <span
                  v-for="(teacher, teacherIndex) in getAllTeachers(examGroup)"
                  :key="teacherIndex"
                  class="px-2 py-1 mb-1 mr-1 text-xs text-indigo-800 bg-indigo-100 rounded-full"
                >
                  {{ teacher.name }}
                </span>
              </div>
            </div>

            <!-- Expand button for more details -->
            <div class="flex justify-center w-full">
              <button
                @click="toggleExamDetails(weekIndex, examIndex)"
                class="flex items-center px-3 py-1 text-sm rounded-full textsm bg-muted"
              >
                {{ examGroup.showDetails ? "Hide details" : "Show details" }}
                <svg
                  class="w-4 h-4 transition-transform duration-200"
                  :class="{ 'transform rotate-180': examGroup.showDetails }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 9l-7 7-7-7"
                  ></path>
                </svg>
              </button>
            </div>

            <!-- Expanded details -->
            <div
              v-if="examGroup.showDetails"
              class="pt-3 mt-3 text-sm border-t border-gray-200"
            >
              <!-- Multiple exams info -->
              <div v-if="examGroup.exams.length > 1">
                <h5 class="mb-2 font-medium">Combined Exams:</h5>
                <div
                  v-for="(exam, examIdx) in examGroup.exams"
                  :key="examIdx"
                  class="pl-2 mb-3 border-l-2 border-gray-200"
                >
                  <div class="flex justify-between">
                    <span class="font-medium capitalize">{{
                      exam.subject
                    }}</span>
                    <span
                      :class="{
                        'bg-lime-300 text-lime-900': exam.complete,
                        'bg-orange-300 text-orange-900': !exam.complete,
                      }"
                      class="px-2 py-1 text-xs font-medium rounded-full"
                    >
                      {{ exam.complete ? "Complete" : "Incomplete" }}
                    </span>
                  </div>
                  <p class="text-gray-600">Duration: {{ exam.duration }} min</p>
                  <p class="text-gray-600">
                    Time: {{ formatTime(exam.start) }} -
                    {{ formatTime(exam.end) }}
                  </p>

                  <!-- Teacher info per exam -->
                  <div
                    v-if="exam.bookedSegments && exam.bookedSegments.length > 0"
                    class="mt-1"
                  >
                    <div
                      v-for="(segment, segmentIdx) in exam.bookedSegments"
                      :key="segmentIdx"
                      class="text-gray-600"
                    >
                      <p>Teacher: {{ segment.teacher.name }}</p>
                      <p>Subjects: {{ segment.teacher.subjects.join(", ") }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Single exam teacher details -->
              <div v-else>
                <div
                  v-if="
                    examGroup.exams[0].bookedSegments &&
                    examGroup.exams[0].bookedSegments.length > 0
                  "
                >
                  <div
                    v-for="(segment, segmentIndex) in examGroup.exams[0]
                      .bookedSegments"
                    :key="segmentIndex"
                    class="mb-2"
                  >
                    <h5 class="font-medium">
                      Teacher: {{ segment.teacher.name }}
                    </h5>
                    <p class="">
                      Subjects: {{ segment.teacher.subjects.join(", ") }}
                    </p>
                    <p class="">Availability on this day:</p>
                    <ul class="pl-5 list-disc">
                      <li
                        v-for="(avail, availIndex) in filterAvailabilitiesByDow(
                          segment.teacher.availabilities,
                          examGroup.exams[0].dow
                        )"
                        :key="availIndex"
                      >
                        {{ formatTime(avail.start) }} -
                        {{ formatTime(avail.end) }}
                      </li>
                    </ul>
                  </div>
                </div>
                <div v-else class="text-gray-600">No teacher assigned yet.</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";

// Define interfaces based on your existing types
interface Teacher {
  id: string;
  name: string;
  availabilities: {
    dow: number;
    start: string | Date;
    end: string | Date;
  }[];
  subjects: string[];
  bias: number;
}

interface BookedSegment {
  teacher: Teacher;
  start: string | Date;
  end: string | Date;
}

interface Exam {
  id: string;
  subject: string;
  start: string | Date;
  end: string | Date;
  duration: number;
  room: string;
  bookedSegments: BookedSegment[];
  complete: boolean;
  dow: number;
}

// Props and emits
const props = defineProps<{
  exams: Exam[];
  loading?: boolean;
}>();

// State
const localLoading = ref(props.loading ?? false);

// Watch for prop changes
watch(
  () => props.loading,
  (newVal) => {
    if (newVal !== undefined) {
      localLoading.value = newVal;
    }
  }
);

// Interface for grouped exams (merged by room and start time)
interface ExamGroup {
  exams: Exam[];
  room: string;
  start: string | Date;
  showDetails: boolean;
}

// Interface for weekly grouped exams
interface WeekGroup {
  weekStart: Date;
  weekEnd: Date;
  weekLabel: string;
  examGroups: ExamGroup[];
}

// Group exams by room and start time
const groupExamsByRoomAndTime = (exams: Exam[]): ExamGroup[] => {
  const groups: Record<string, ExamGroup> = {};

  exams.forEach((exam) => {
    // Create a unique key for each room and start time combination
    const startString =
      typeof exam.start === "string" ? exam.start : exam.start.toISOString();
    const key = `${exam.room}-${startString}`;

    if (!groups[key]) {
      groups[key] = {
        exams: [],
        room: exam.room,
        start: exam.start,
        showDetails: false,
      };
    }

    groups[key].exams.push({ ...exam });
  });

  return Object.values(groups);
};

const localWeeklyExams = ref<WeekGroup[]>([]);

watch(
  () => props.exams,
  (newExams) => {
    if (!newExams || newExams.length === 0) {
      localWeeklyExams.value = [];
      return;
    }

    // Group by week and group by room+time
    const weeks: Record<string, WeekGroup> = {};

    newExams.forEach((exam) => {
      const examDate = new Date(exam.start);
      const weekStart = new Date(examDate);
      weekStart.setDate(
        examDate.getDate() -
          examDate.getDay() +
          (examDate.getDay() === 0 ? -6 : 1)
      );

      const weekKey = weekStart.toISOString().split("T")[0];

      if (!weeks[weekKey]) {
        const weekEnd = new Date(weekStart);
        weekEnd.setDate(weekStart.getDate() + 6);

        weeks[weekKey] = {
          weekStart,
          weekEnd,
          weekLabel: formatWeekLabel(weekStart, weekEnd),
          examGroups: [],
        };
      }
    });

    Object.keys(weeks).forEach((weekKey) => {
      const weekStart = weeks[weekKey].weekStart;
      const weekEnd = weeks[weekKey].weekEnd;

      const weekExams = newExams.filter((exam) => {
        const examDate = new Date(exam.start);
        return examDate >= weekStart && examDate <= weekEnd;
      });

      weeks[weekKey].examGroups = groupExamsByRoomAndTime(weekExams);
    });

    // Finally assign
    localWeeklyExams.value = Object.values(weeks).sort(
      (a, b) => a.weekStart.getTime() - b.weekStart.getTime()
    );
  },
  { immediate: true }
);

// Group exams by week
// const weeklyExams = computed<WeekGroup[]>(() => {
//   if (!props.exams || props.exams.length === 0) return [];

//   // Group exams by week
//   const weeks: Record<string, WeekGroup> = {};

//   props.exams.forEach((exam) => {
//     const examDate = new Date(exam.start);
//     // Get week start (Monday)
//     const weekStart = new Date(examDate);
//     weekStart.setDate(
//       examDate.getDate() -
//         examDate.getDay() +
//         (examDate.getDay() === 0 ? -6 : 1)
//     );

//     // Format week key: YYYY-MM-DD
//     const weekKey = weekStart.toISOString().split("T")[0];

//     // Create week if it doesn't exist
//     if (!weeks[weekKey]) {
//       const weekEnd = new Date(weekStart);
//       weekEnd.setDate(weekStart.getDate() + 6);

//       weeks[weekKey] = {
//         weekStart,
//         weekEnd,
//         weekLabel: formatWeekLabel(weekStart, weekEnd),
//         examGroups: [],
//       };
//     }
//   });

//   // For each week, group exams by room and start time
//   Object.keys(weeks).forEach((weekKey) => {
//     const weekStart = weeks[weekKey].weekStart;
//     const weekEnd = weeks[weekKey].weekEnd;

//     // Filter exams that fall within this week
//     const weekExams = props.exams.filter((exam) => {
//       const examDate = new Date(exam.start);
//       return examDate >= weekStart && examDate <= weekEnd;
//     });

//     // Group these exams by room and start time
//     weeks[weekKey].examGroups = groupExamsByRoomAndTime(weekExams);
//   });

//   // Convert to array and sort by week start date
//   return Object.values(weeks).sort(
//     (a, b) => a.weekStart.getTime() - b.weekStart.getTime()
//   );
// });

// Methods
// Format date (Wed, May 8)
const formatDate = (dateString: string | Date): string => {
  const date = new Date(dateString);
  const options: Intl.DateTimeFormatOptions = {
    weekday: "short",
    month: "short",
    day: "numeric",
  };
  return date.toLocaleDateString("en-US", options);
};

// Format time (7:00 AM)
const formatTime = (dateString: string | Date): string => {
  const date = new Date(dateString);
  return date.toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  });
};

// Format week label (May 6 - 12, 2024)
const formatWeekLabel = (startDate: Date, endDate: Date): string => {
  const startMonth = startDate.toLocaleDateString("en-US", { month: "short" });
  const endMonth = endDate.toLocaleDateString("en-US", { month: "short" });

  const startDay = startDate.getDate();
  const endDay = endDate.getDate();

  if (startMonth === endMonth) {
    return `${startMonth} ${startDay} - ${endDay}, ${startDate.getFullYear()}`;
  } else {
    return `${startMonth} ${startDay} - ${endMonth} ${endDay}, ${startDate.getFullYear()}`;
  }
};

// Filter availabilities by day of week
const filterAvailabilitiesByDow = (availabilities: any[], dow: number) => {
  return availabilities.filter((a) => a.dow === dow);
};

// Toggle exam details
const toggleExamDetails = (weekIndex: number, examIndex: number) => {
  localWeeklyExams.value[weekIndex].examGroups[examIndex].showDetails =
    !localWeeklyExams.value[weekIndex].examGroups[examIndex].showDetails;
  console.log(
    "Toggle show, week: ",
    weekIndex,
    examIndex,
    "State: ",
    localWeeklyExams.value[weekIndex].examGroups[examIndex].showDetails
  );
};

// Get border color class based on exam completion status
const getBorderClass = (examGroup: ExamGroup): string => {
  const isSingleExam = examGroup.exams.length === 1;
  const allComplete = examGroup.exams.every((exam) => exam.complete);
  const allIncomplete = examGroup.exams.every((exam) => !exam.complete);

  if (isSingleExam) {
    return examGroup.exams[0].complete
      ? "border-transparent" // Single exam complete = Success
      : "border-red-500"; // Single exam incomplete = Warning
  }

  if (allComplete) {
    return "border-transparent"; // All exams complete = Success
  }

  if (allIncomplete) {
    return "border-red-500"; // All exams incomplete = Warning
  }

  return "border-orange-500"; // Mixed completion = Info
};

// Get combined subjects label
const getSubjectsLabel = (examGroup: ExamGroup): string => {
  if (examGroup.exams.length === 1) {
    return examGroup.exams[0].subject;
  }

  // Multiple subjects, limit to 2 with "and more" text
  const subjects = [...new Set(examGroup.exams.map((exam) => exam.subject))];

  if (subjects.length <= 2) {
    return subjects.join(" & ");
  }

  return `${subjects[0]} & ${subjects.length - 1} more`;
};

// Get latest end time from all exams in a group
const getLatestEndTime = (examGroup: ExamGroup): string => {
  const endTimes = examGroup.exams.map((exam) => new Date(exam.end).getTime());
  const latestTime = new Date(Math.max(...endTimes));
  return latestTime.toISOString();
};

// Get duration range (e.g., "30-45 min" or just "30 min" if all same)
const getDurationRange = (examGroup: ExamGroup): string => {
  const durations = examGroup.exams.map((exam) => exam.duration);
  const minDuration = Math.min(...durations);
  const maxDuration = Math.max(...durations);

  if (minDuration === maxDuration) {
    return `${minDuration} min`;
  }

  return `${minDuration}-${maxDuration} min`;
};

// Get all teachers from all exams in a group
const getAllTeachers = (examGroup: ExamGroup): Teacher[] => {
  const teachers: Teacher[] = [];
  const teacherIds = new Set<string>();

  examGroup.exams.forEach((exam) => {
    if (exam.bookedSegments) {
      exam.bookedSegments.forEach((segment) => {
        if (!teacherIds.has(segment.teacher.id)) {
          teacherIds.add(segment.teacher.id);
          teachers.push(segment.teacher);
        }
      });
    }
  });

  return teachers;
};
</script>
