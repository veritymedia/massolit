<template>
  <div>
    <ul class="flex gap-4">
      <Card
        v-for="t in timetables"
        :key="t.id"
        class="flex flex-col gap-2 p-4 w-72 bg-background"
      >
        <h3 class="font-bold">{{ t.name }}</h3>
        <div class="flex items-baseline justify-between">
          Teachers
          <Button
            @click="navigateTo(`/timetables/${t.id}/teachers`)"
            variant="ghost"
            >Edit</Button
          >
        </div>
        <div class="flex items-baseline justify-between">
          Exams to plan
          <Button
            @click="navigateTo(`/timetables/${t.id}/exams`)"
            variant="ghost"
            >Edit</Button
          >
        </div>
        <div v-if="t.timetable.length > 0" class="flex gap-2">
          <Button @click="navigateTo(`/timetables/${t.id}`)" class="w-full"
            >View</Button
          >
          <DangerButton
            @first-click="console.log('confirm_to_recalculate')"
            @confirmed="calculateTimetable(t.id)"
            variant="secondary"
            :show-countdown="false"
          >
            <template #default-content>
              <Icon
                class="size-5"
                name="material-symbols:rotate-left-rounded"
              />
            </template>
            <template #confirm-content>
              <Icon
                class="size-5"
                name="material-symbols:check-small-rounded"
              />
            </template>
          </DangerButton>
          <DangerButton
            variant="destructive"
            @first-click="console.log('confirm_to_delete')"
            @confirmed="deleteTimetable(t.id)"
            :show-countdown="false"
          >
            <template #default-content>
              <Icon class="size-5" name="material-symbols:delete-outline" />
            </template>
            <template #confirm-content>
              <Icon
                class="size-5"
                name="material-symbols:check-small-rounded"
              />
            </template>
          </DangerButton>
        </div>
        <div class="flex w-full gap-2" v-else>
          <Button
            @click="calculateTimetable(t.id)"
            variant="secondary"
            class="w-full"
            >Calculate</Button
          >
          <Button variant="destructive" @click="deleteTimetable(t.id)"
            ><Icon name="material-symbols:delete-outline"
          /></Button>
        </div>
      </Card>
    </ul>
  </div>
</template>
<script setup lang="ts">
import { type TeacherRaw, Teacher, Exam, type ExamRaw } from "timetabler";
import { DangerButton } from "../ui/danger-button";
import {
  assignTeachersToExams,
  type Configuration,
} from "~/timetabler/core/v2";
const pb = usePocketbase();

interface Props {
  newTimetable: Timetable;
}
const props = defineProps<Props>();

export type Timetable = {
  id: string;
  name: string;
  exam_timetable: [];
  teachers: [];
  timetable: [];
};

watch(
  () => props.newTimetable,
  (curr, old) => {
    console.log("watch: new time", curr, old);
    if (!old && !curr) {
      return;
    }
    if (curr.id !== old.id) {
      timetables.value.unshift(curr);
    }
  },
);

async function calculateTimetable(id: string) {
  try {
    const timetable = await pb.collection("timetables").getOne(id, {
      expand: "teachers,exam_timetables",
    });
    let examsRaw: ExamRaw[] = [];

    timetable.expand.exam_timetables.forEach((list: any) => {
      examsRaw = examsRaw.concat(list.data);
    });

    // const exams = examsRaw.map((exam) => {
    //   return new Exam(exam);
    // });
    // console.log("TM: ", timetable.expand.exam_timetables);
    console.log("TIMETABLE: ", timetable.expand.teachers);

    const teachers: Teacher[] = timetable.expand.teachers.map(
      (teacher: TeacherRaw) => {
        return new Teacher({
          name: teacher.name,
          availabilities: teacher.availabilities,
          subjects: teacher.subjects,
        });
      },
    );

    // const finalTimetable = processExams(exams, teachers);
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

    const finalTimetableV2 = assignTeachersToExams(
      examsRaw,
      timetable.expand.teachers,
      config,
    );
    console.log("FINAL: ", finalTimetableV2);

    const saved = await pb.collection("timetables").update(id, {
      timetable: finalTimetableV2,
    });
    console.log(saved);
    return navigateTo(`/timetables/${id}`);
  } catch (error) {
    console.log(error);
  }
}

async function getTimetables() {
  try {
    const records = await pb
      .collection("timetables")
      .getFullList<Timetable>(1, {
        sort: "-created",
      });
    timetables.value = records;
    console.log(records);
  } catch (err) {
    console.log(err);
  }
}

async function deleteTimetable(id: string) {
  try {
    await pb.collection("timetables").delete(id);
  } catch (err) {
    console.log(err);
  } finally {
    getTimetables();
  }
}

onMounted(async () => {
  await getTimetables();
});

const timetables = ref();
</script>
