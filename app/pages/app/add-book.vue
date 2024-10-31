<template>
  <div class="flex flex-col gap-2 overflow-scroll">
    <div class="flex flex-col items-center">
      <Icon
        name="material-symbols:warning-rounded"
        class="w-12 h-12 text-[yellow] opacity-75"
      />
      <p class="w-full text-center my-5">
        This book does not exist in your database. <br />
      </p>
    </div>
    <div v-if="foundBookList" class="">
      <Card v-for="book in foundBookList" :key="book.id" class="flex p-3">
        <div>
          <h2>{{ book.volumeInfo.title }}</h2>
          <div class="flex gap-2 uppercase text-xs opacity-50">
            {{ reduceAuthorsToString(book.volumeInfo.authors) }}
          </div>
        </div>
        <img
          class="max-w-32 place-self-center"
          v-if="getBestBookCoverURL(book.volumeInfo.imageLinks)"
          :src="getBestBookCoverURL(book.volumeInfo.imageLinks)"
        />
      </Card>
    </div>
    <Button
      v-if="foundBookList"
      @click="addBookToDatabase"
      class="w-full"
      variant="secondary"
      >Add book</Button
    >
    <div v-if="bookNotFound">
      <Card class="p-3 flex flex-col gap-3">
        <h2>Book could not be found. Add it manually.</h2>
        <div>
          <label for="bookTitle">
            Title <span class="font-bold text-destructive">*</span></label
          >
          <Input v-model="manualBook.title" id="bookTitle" />
        </div>
        <div>
          <label for="bookCoverUrl"> Cover URL</label>
          <Input id="bookCoverUrl" v-model="manualBook.coverUrl" />
        </div>
        <Button class="w-full" variant="secondary" @click="addBookToDatabase"
          >Add book</Button
        >
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { RentedBookStatus as BookStatus } from "../app.vue";
import { useAlert } from "#imports";

interface Props {
  rentedBookStatus: BookStatus;
}
const props = withDefaults(defineProps<Props>(), {});

const { createAlert } = useAlert();

const pb = usePocketbase();

type IndustryIdentifier = {
  type: string;
  identifier: string;
};
function getIsbn13Identifier(identifiers: IndustryIdentifier[]): string | null {
  const isbn13 = identifiers.find((id) => id.type === "ISBN_13");
  return isbn13 ? isbn13.identifier : null;
}

const manualBook = ref({
  title: "",
  coverUrl: "",
});

async function addBookToDatabase() {
  try {
    if (bookNotFound.value) {
      if (!manualBook.value.title || manualBook.value.title.length < 5) {
        createAlert({
          title: "Book title issue",
          message: "Book title too short. Min length is 5 characters",
        });
        console.log("Book title too short: ", manualBook.value.title);
      }
      const data = {
        title: manualBook.value.title,
        cover_url: manualBook.value.coverUrl ?? "",
        isbn: props.rentedBookStatus.parsedCode.isbn,
      };

      const record = await pb.collection("books").create(data);
      console.log(record);
      await navigateTo(
        `/app?code=${props.rentedBookStatus.parsedCode.rawCode}`,
      );
      return;
    } else {
      const volumeInfo = foundBookList.value[0].volumeInfo;
      console.log("resource", volumeInfo);
      const isbn = getIsbn13Identifier(volumeInfo.industryIdentifiers);

      if (isbn === null) {
        throw new Error("no available ISBN 13 in retrieved book.");
      }

      const data = {
        title: volumeInfo.title,
        isbn: props.rentedBookStatus.parsedCode.isbn,
        cover_url: getBestBookCoverURL(volumeInfo.imageLinks),
      };
      const record = await pb.collection("books").create(data);
      console.log(record);
    }
  } catch (err) {}
}

function reduceAuthorsToString(authors: string[]): string {
  return authors.reduce((v, str) => {
    return (str = `${str}, ${v}`);
  });
}

type ImageLinks = {
  smallThumbnail?: string;
  thumbnail?: string;
  small?: string;
  medium?: string;
  large?: string;
  extraLarge?: string;
};

function getBestBookCoverURL(images: ImageLinks): string {
  if (images.extraLarge) return images.extraLarge;
  if (images.large) return images.large;
  if (images.medium) return images.medium;
  if (images.small) return images.small;
  if (images.thumbnail) return images.thumbnail;
  if (images.smallThumbnail) return images.smallThumbnail;

  return "";
}

const foundBookList = ref();

async function searchBookByISBN(
  isbn: string | undefined,
): Promise<object | undefined> {
  try {
    console.log("isbn", isbn);
    const res = await fetch(
      `https://www.googleapis.com/books/v1/volumes?q=isbn:${isbn}`,
      {},
    );
    const data = await res.json();
    console.log("Books by ISBN: ", data);
    if (data.totalItems === 0) {
      console.log("No books found for isbn ", isbn);
      createAlert({
        title: "Book not found.",
        message: "Could not find book by ISBN. Please search for it by title.",
      });
      return undefined;
    }
    return data.items[0];
  } catch (err) {
    console.log(err);
  }
}

const bookNotFound = ref(false);

onMounted(async () => {
  console.log("prop: ", props.rentedBookStatus);

  const book = await searchBookByISBN(props.rentedBookStatus.parsedCode.isbn);

  if (book !== undefined) {
    foundBookList.value = [book];
  } else {
    bookNotFound.value = true;
  }
});
</script>
