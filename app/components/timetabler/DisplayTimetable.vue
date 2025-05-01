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
          <div class="w-3 h-3 mr-2 bg-purple-500 rounded-full"></div>
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
            <!-- Subject(s) and exam code badge -->
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
                v-if="examGroup.exams[0].examCode"
                class="px-2 py-1 text-xs font-medium text-blue-800 bg-blue-100 rounded-full"
              >
                {{ examGroup.exams[0].examCode }}
              </span>
              <span
                v-else-if="examGroup.exams.length > 1"
                class="px-2 py-1 text-xs font-medium text-blue-800 bg-blue-100 rounded-full"
              >
                Multiple
              </span>
            </div>

            <!-- Time and room info -->
            <div class="flex flex-col gap-2">
              <div>
                <span class="text-lg font-bold">{{
                  formatTimeFromString(examGroup.exams[0].start)
                }}</span>
              </div>
              <div class="flex flex-col gap-2">
                <div>
                  <span class="pr-2 text-xl">📅</span>
                  {{ formatDateFromString(examGroup.exams[0].date) }}
                </div>
                <div>
                  <span class="pr-2 text-xl">⏳</span>
                  {{ getEndTime(examGroup.exams[0]) }} ({{
                    formatDuration(examGroup.exams[0].duration)
                  }})
                </div>

                <div>
                  <span class="pr-2 text-xl">🏢</span>
                  {{ examGroup.room }}
                </div>
              </div>
            </div>

            <!-- Proctors -->
            <div v-if="getAllProctors(examGroup).length > 0" class="mt-3">
              <h4 class="mb-1 text-sm font-medium">Proctored by:</h4>
              <div class="flex flex-wrap">
                <span
                  v-for="(proctor, proctorIndex) in getAllProctors(examGroup)"
                  :key="proctorIndex"
                  class="px-2 py-1 mb-1 mr-1 text-xs text-indigo-800 bg-indigo-100 rounded-full"
                >
                  {{ proctor }}
                </span>
              </div>
            </div>

            <!-- Expand button for more details -->
            <div class="flex justify-center w-full mt-3">
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
                      v-if="exam.examCode"
                      class="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-900 rounded-full"
                    >
                      {{ exam.examCode }}
                    </span>
                  </div>
                  <p class="text-gray-600">
                    Duration: {{ formatDuration(exam.duration) }}
                  </p>
                  <p class="text-gray-600">
                    Time: {{ formatTimeFromString(exam.start) }} -
                    {{ getEndTime(exam) }}
                  </p>

                  <!-- Proctor info per exam -->
                  <div
                    v-if="exam.proctors && exam.proctors.length > 0"
                    class="mt-1"
                  >
                    <p class="font-medium">Proctors:</p>
                    <div
                      v-for="(proctor, proctorIdx) in exam.proctors"
                      :key="proctorIdx"
                      class="text-gray-600 pl-2"
                    >
                      <p>
                        {{ proctor.teacher }} ({{ proctor.start }} -
                        {{ proctor.end }})
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Single exam proctor details -->
              <div v-else>
                <div
                  v-if="
                    examGroup.exams[0].proctors &&
                    examGroup.exams[0].proctors.length > 0
                  "
                >
                  <h5 class="font-medium mb-2">Proctor Schedule:</h5>
                  <div
                    v-for="(proctor, proctorIndex) in examGroup.exams[0]
                      .proctors"
                    :key="proctorIndex"
                    class="mb-2 pl-2 border-l-2 border-gray-200"
                  >
                    <p class="font-medium">{{ proctor.teacher }}</p>
                    <p class="text-gray-600">
                      {{ proctor.start }} - {{ proctor.end }}
                    </p>
                  </div>
                </div>
                <div v-else class="text-gray-600">
                  No proctors assigned yet.
                </div>
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

// Define interfaces based on the new data structure
interface Proctor {
  teacher: string;
  start: string;
  end: string;
}

interface Exam {
  id: string;
  subject: string;
  date: string;
  start: string;
  duration: string;
  room: string;
  examCode?: string;
  proctors: Proctor[];
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
  },
);

// Interface for grouped exams (merged by room and start time)
interface ExamGroup {
  exams: Exam[];
  room: string;
  date: string;
  start: string;
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
    // Create a unique key for each room, date and start time combination
    const key = `${exam.room}-${exam.date}-${exam.start}`;

    if (!groups[key]) {
      groups[key] = {
        exams: [],
        room: exam.room,
        date: exam.date,
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
      const examDate = new Date(exam.date);
      const weekStart = new Date(examDate);
      weekStart.setDate(
        examDate.getDate() -
          examDate.getDay() +
          (examDate.getDay() === 0 ? -6 : 1),
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
        const examDate = new Date(exam.date);
        return examDate >= weekStart && examDate <= weekEnd;
      });

      weeks[weekKey].examGroups = groupExamsByRoomAndTime(weekExams);
    });

    // Finally assign
    localWeeklyExams.value = Object.values(weeks).sort(
      (a, b) => a.weekStart.getTime() - b.weekStart.getTime(),
    );
  },
  { immediate: true },
);

// Format date from string (Wed, May 8)
const formatDateFromString = (dateString: string): string => {
  const date = new Date(dateString);
  const options: Intl.DateTimeFormatOptions = {
    weekday: "short",
    month: "short",
    day: "numeric",
  };
  return date.toLocaleDateString("en-US", options);
};

// Format time from string (7:00 AM)
const formatTimeFromString = (timeString: string): string => {
  // Handle ISO format or just time format
  const date = timeString.includes("T")
    ? new Date(timeString)
    : new Date(`2000-01-01T${timeString}`);

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

// Format duration from string "HH:MM" to "XX min"
const formatDuration = (duration: string): string => {
  if (!duration) return "";

  // If already in minutes format
  if (duration.toString().indexOf(":") === -1) {
    return `${duration} min`;
  }

  // Convert HH:MM to minutes
  const [hours, minutes] = duration.split(":").map(Number);
  const totalMinutes = hours * 60 + minutes;
  return `${totalMinutes} min`;
};

// Calculate end time based on start time and duration
const getEndTime = (exam: Exam): string => {
  // Parse start time
  let startDate: Date;

  if (exam.start.includes("T")) {
    // ISO format
    startDate = new Date(exam.start);
  } else {
    // Time only format, use exam date
    startDate = new Date(`${exam.date}T${exam.start}`);
  }

  // Parse duration and add to start time
  const [hours, minutes] = exam.duration.split(":").map(Number);
  const durationMs = (hours * 60 + minutes) * 60 * 1000;

  const endDate = new Date(startDate.getTime() + durationMs);

  return endDate.toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  });
};

// Toggle exam details
const toggleExamDetails = (weekIndex: number, examIndex: number) => {
  localWeeklyExams.value[weekIndex].examGroups[examIndex].showDetails =
    !localWeeklyExams.value[weekIndex].examGroups[examIndex].showDetails;
};

// Get border color class based on exam properties
const getBorderClass = (examGroup: ExamGroup): string => {
  // In the new data structure, we don't have a 'complete' field
  // You can implement your own logic based on other fields
  // For now, using a fixed class
  return "border-transparent";
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

// Get all proctors from all exams in a group
const getAllProctors = (examGroup: ExamGroup): string[] => {
  const proctorSet = new Set<string>();

  examGroup.exams.forEach((exam) => {
    if (exam.proctors) {
      exam.proctors.forEach((proctor) => {
        proctorSet.add(proctor.teacher);
      });
    }
  });

  return Array.from(proctorSet);
};
</script>
