export const useDialogState = () => {
  const isOpen = ref(false);

  function closeDialog() {
    isOpen.value = false;
  }

  function openDialog() {
    isOpen.value = true;
  }

  return [isOpen, closeDialog(), openDialog()];
};
