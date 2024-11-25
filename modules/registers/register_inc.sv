module register_inc #(
    parameter N = 2
) (
    input  logic [N-1:0] data_in,
    input  logic         write, read, increment,
    output logic [N-1:0] data_out,
);

    logic [N-1:0] data;

    // always_ff @(posedge write, posedge increment) begin
    //     if (write)
    //         data <= data_in;
    //     else if (increment)
    //         data <= data + 1;
    // end
    always_latch  begin
        if (write)
            data <= data_in;
        else if (increment)
            data <= data + 1;
    end

    always_comb begin
        if (read)
            data_out <= data;
        else
            data_out <= 0;
    end

endmodule