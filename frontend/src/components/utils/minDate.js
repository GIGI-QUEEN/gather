const minDate = (now) => {
    const year = now.getFullYear();
    const month = ("0" + (now.getMonth() + 1)).slice(-2);
    const day = ("0" + now.getDate()).slice(-2);
    const hour = ("0" + now.getHours()).slice(-2);
    const minute = ("0" + now.getMinutes()).slice(-2);
    const minimumDate =
      year + "-" + month + "-" + day + "T" + hour + ":" + minute;
    return minimumDate;
  };

  export default minDate