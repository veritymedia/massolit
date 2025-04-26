<script lang="ts" setup>
import type { Timetable } from "~/components/timetabler/TimetableList.vue";
import TimetableList from "~/components/timetabler/TimetableList.vue";
// const isLoading = ref(false);

const pb = usePocketbase();

async function createTimetable() {
  try {
    const now = new Date();
    const data = {
      name: `Timetable ${now.toDateString()}`,
      exam_timetable: "",
      teachers: [],
      timetable: "[]",
    };

    newTimetable.value = await pb.collection("timetables").create(data);
  } catch (err) {
    console.log(err);
  }
}

const newTimetable = ref({
  id: "",
  name: "",
  exam_timetable: [],
  teachers: [],
  timetable: [],
});
</script>

<template>
  <div class="">
    <div
      v-if="false"
      class="fixed top-0 left-0 flex items-center justify-center w-screen h-screen bg-background"
    >
      <Icon name="line-md:loading-loop" class="w-16 h-16" />
    </div>
    <AppAlert />
    <div>
      <div class="flex items-baseline gap-2">
        <h2 class="mt-10 mb-5">Exam Timetables</h2>
        <Button @click="createTimetable">New Timetable</Button>
      </div>
      <div><TimetableList :newTimetable="newTimetable" /></div>
    </div>
  </div>
</template>

<style></style>
