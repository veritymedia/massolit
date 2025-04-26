<script lang="ts" setup>
import DisplayTimetable from "~/components/timetabler/DisplayTimetable.vue";
const pb = usePocketbase();
const route = useRoute();

const timetableState = ref();
const isLoading = ref(true);

async function getTimetable() {
  const id = route.params.timetableId as string;

  try {
    const tb = await pb.collection("timetables").getOne(id);
    console.log(tb);
    timetableState.value = tb.timetable;
    isLoading.value = false;
  } catch (error) {
    console.log(error);
  }
}

onMounted(async () => {
  await getTimetable();
});
</script>

<template>
  <div class="flex flex-col gap-5 mt-10 bg-white">
    <div class="flex gap-2">
      <NuxtLink to="/timetables"
        ><Icon class="size-6" name="material-symbols:arrow-left-alt-rounded"
      /></NuxtLink>
      <h2>Timetable</h2>
    </div>
    <DisplayTimetable :loading="isLoading" :exams="timetableState" />
  </div>
</template>

<style></style>
