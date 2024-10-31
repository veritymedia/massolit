<template>
  <Card class="p-3 flex flex-col items-start gap-2">
    <h2 class="text-lg font-bold">
      {{ props.book?.title }}
    </h2>

    <div class="flex text-xs justify-center gap-2 items-start">
      <p class="bg-accent rounded-full px-2 py-0.5">
        {{ props.book?.isbn }}
      </p>
      <p v-if="!props.bookInstanceMissing" class="rounded-full px-2 py-0.5">
        {{ props.book.id }}
      </p>
    </div>
    <img
      v-if="props.book.cover_url"
      class="w-1/2 place-self-center mt-4"
      :src="props.book.cover_url"
      :alt="'book-cover-' + props.book.title"
    />
    <div
      v-else
      class="aspect-[9/16] bg-muted flex-col p-3 h-32 place-self-center flex items-center text-center rounded-lg justify-center mt-4"
    >
      <Icon
        name="material-symbols:question-mark-rounded"
        class="w-12 h-12 text-foreground bg-opacity-75"
      />
      <p>No cover</p>
    </div>
  </Card>
</template>

<script setup lang="ts">
interface Props {
  book: {
    id: string;
    isbn: string;
    title: string;
    cover_url: string;
  };
  bookInstanceMissing?: boolean;
}

const props = defineProps<Props>();
</script>
