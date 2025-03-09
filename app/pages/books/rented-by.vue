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
    <Dialog v-model:open="bookReturnModal" class="">
      <DialogTrigger class="sticky bottom-5" as-child>
        <Button class="w-full" variant="secondary">Return Book</Button>
      </DialogTrigger>
      <DialogContent class="h-full flex flex-col items-center justify-between">
        <DialogHeader>
          <DialogTitle>Return book?</DialogTitle>
        </DialogHeader>

        <div class="w-full flex flex-col items-center gap-4">
          <h2 class="text-lg font-bold"></h2>
          <p class="text-center">
            The book
            <span class="font-bold">
              {{ props.rentedBookStatus.book?.title }}
            </span>
            will be removed from account of
            <span class="font-bold">
              {{ user?.first_name + " " + user?.last_name }} </span
            >.
          </p>
          <Button variant="destructive" @click="handleBookReturn"
            >Return Book</Button
          >
          <Button variant="ghost" @click="closeDialog">Cancel</Button>
        </div>
        <DialogFooter class="w-full flex items-center"> </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import type { Record } from "pocketbase";
import type { RentedBookStatus as BookStatus } from "../app.vue";

interface Props {
  rentedBookStatus: BookStatus;
}
const props = withDefaults(defineProps<Props>(), {});

const pb = usePocketbase();

const { isOpen: bookReturnModal, closeDialog, openDialog } = useDialogState();

async function handleBookReturn() {
  try {
    const rental = await getRentalRecord();

    if (rental === undefined) {
      return;
    }

    const recordDeleted = await pb.collection("rentals").delete(rental.id);
    if (recordDeleted) {
      await navigateTo("/books");
    }
  } catch (err) {
    console.log(err);
  }
}

async function getRentalRecord(): Promise<Record | undefined> {
  try {
    const filter = `(book_instance="${props.rentedBookStatus.rental?.book_instance_id}"&&rented_to="${props.rentedBookStatus.rental?.managebac_user_id}")`;

    const record = await pb.collection("rentals").getFirstListItem(filter);

    return record;
  } catch (err) {
    console.log(err);
    return undefined;
  }
}

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
    !props.rentedBookStatus.rental?.managebac_user_id ||
    !props.rentedBookStatus.bookId ||
    !props.rentedBookStatus.book?.isbn
  ) {
    await navigateTo("/books");
  }
  if (props.rentedBookStatus.rental?.managebac_user_id) {
    await getManagebacUser(props.rentedBookStatus.rental.managebac_user_id);
  }
});
</script>
