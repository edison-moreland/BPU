module is_control_char (
    input logic [7:0] char,
    output logic is_printable,
    output logic is_control,
);

    always_comb begin
        is_printable = 0;
        is_control = 1;

        if ((char >= 'h20) && (char < 'h7F)) begin
            is_printable = 1;
            is_control = 0;
        end
    end


endmodule