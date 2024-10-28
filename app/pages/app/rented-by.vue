<template>
  <div class="flex flex-col items-center">
    <Icon name="material-symbols:check-circle" class="w-12 h-12 opacity-70" />
    <p v-if="user" class="w-full text-center my-5">
      Book borrowed by <br />
      <span class="text-lg uppercase font-bold">{{
        user.first_name + " " + user.last_name
      }}</span>
      <br />
      {{ user.class_grade }}
    </p>
  </div>
  <BookCard
    class="my-5"
    :book="{
      id: props.rentedBookStatus.bookId,
      isbn: props.rentedBookStatus.book?.isbn,
      title: props.rentedBookStatus.book?.title,
      cover_url: props.rentedBookStatus.book?.cover_url,
    }"
  />
  <div class="w-full">
    <Button class="w-full" variant="secondary">Return Book</Button>
  </div>
</template>

<script setup lang="ts">
import type { RentedBookStatus as BookStatus } from "../app.vue";

interface Props {
  rentedBookStatus: BookStatus;
}
const props = withDefaults(defineProps<Props>(), {});

const pb = usePocketbase();

type ManagebacStudent = {
  account_uid?: string;
  archived?: boolean;
  class_grade?: string;
  class_grade_number?: number;
  created_at?: string; // ISO date string
  email: string;
  first_name: string;
  middle_name: string;
  graduating_year?: number;
  homeroom_advisor_id?: number;
  id: number;
  languages?: string[];
  last_accessed_at?: string; // ISO date string
  last_name: string;
  nationalities?: string[];
  parent_ids?: number[];
  program?: string;
  program_code?: string;
  role?: string;
  timezone?: string;
  ui_language?: string;
  updated_at?: string; // ISO date string
  year_group_id?: number;
};

const user = ref<ManagebacStudent>();

async function getManagebacUser(managebacID: string) {
  console.log("Fetching ManageBac Student with id: ", managebacID);

  try {
    const res = await pb.send(`/managebac/students/${managebacID}`, {});
    user.value = res.student;
    console.log(res);
  } catch (error) {
    console.log(error);
  }
}

onMounted(async () => {
  if (
    !props.rentedBookStatus.renter?.id ||
    !props.rentedBookStatus.bookId ||
    !props.rentedBookStatus.book?.isbn
  ) {
    await navigateTo("/app");
  }
  if (props.rentedBookStatus.renter?.id && user) {
    await getManagebacUser(props.rentedBookStatus.renter.id);
  }
});
</script>
