<template>
  <div class="flex flex-col gap-2 items-center justify-start">
    <Icon
      name="material-symbols:warning-rounded"
      class="w-12 h-12 text-[yellow]"
    />
    <h2 class="text-lg uppercase font-bold">No such code.</h2>
    <p class="text-center">
      Would you like to add the code
      <span class="bg-primary font-bold text-xs px-2 rounded-full">{{
        props.rentedBookStatus.bookId
      }}</span>
      to the database?
    </p>
    <Button
      variant="secondary"
      class="xw-full"
      :disabled="isLoading"
      @click="handleAddBookInstance()"
      >Add to Database</Button
    >
    <BookCard
      class="mt-5"
      :bookInstanceMissing="!props.rentedBookStatus.codeExists"
      :book="{
        id: props.rentedBookStatus.bookId,
        isbn: props.rentedBookStatus.book?.isbn,
        title: props.rentedBookStatus.book?.title,
        cover_url: props.rentedBookStatus.book?.cover_url,
      }"
    />
  </div>
</template>

<script setup lang="ts">
import type { RentedBookStatus as BookStatus } from "../app.vue";

interface Props {
  rentedBookStatus: BookStatus;
}
const props = withDefaults(defineProps<Props>(), {});

const pb = usePocketbase();

const isLoading = ref(false);
async function handleAddBookInstance() {
  try {
    isLoading.value = true;
    if (!props.rentedBookStatus.book || !props.rentedBookStatus.bookId) {
      return;
    }
    const data = {
      book: props.rentedBookStatus.book.id,
      book_code: props.rentedBookStatus.bookId,
    };

    const res = await pb.collection("book_instances").create(data);
    console.log(res);

    await navigateTo("/app/rent-out");
  } catch (err) {
    console.log(err);
  } finally {
    isLoading.value = false;
  }
}
</script>
